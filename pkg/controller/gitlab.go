package controller

import (
	b64 "encoding/base64"
	"os"
	"regexp"
	"strings"

	"vngitSub/pkg/utils"

	"github.com/xanzy/go-gitlab"
)

func changeTagImage(projectID interface{}, transID string, environment string, imageName string, oldTag string, newTag string, blobList []string, botName string, rootPath string, fullImage string, cluster string) string {
	logger := utils.ConfigZap()
	client := createNewGitlabClient()

	branchName := strings.Split(imageName, "/")[len(strings.Split(imageName, "/"))-1] + "-" + newTag
	oldBranchName := strings.Split(imageName, "/")[len(strings.Split(imageName, "/"))-1] + "-" + oldTag

	if environment == "prod" {
		err,res := createNewBranch(projectID, transID, oldBranchName, branchName)
		if err != nil {
			logger.Errorf("[%s] Creating a new branch named [%s]...%s: %v", transID, branchName, res.Status, err)
			return "FAILED"
		}
		logger.Infof("[%s] Creating a new branch named [%s]...%s", transID, branchName, res.Status)
	}

	for blob := 0; blob < len(blobList); blob++ {
		// Check directory existing
		parentPath := strings.Join(strings.Split(blobList[blob], "/")[0:len(strings.Split(blobList[blob],"/"))-1], "/")
		if _, err := os.Stat(parentPath); os.IsNotExist(err) {
			err := os.MkdirAll(parentPath, os.ModePerm)
			if err != nil {
				logger.Errorf("[%s] Creating a new parent path [%s]...FAILED: %v", transID, parentPath, err)
				return "FAILED"
			}

			logger.Debugf("[%s] Creating a new parent path [%s]...OK", transID, parentPath)
		}

		// Download raw blob files
		logger.Infof("[%s] Downloading raw blob file containing old tag [%s]", transID, oldTag)
		blobRawContent, res, err := client.RepositoryFiles.GetRawFile(projectID, blobList[blob], &gitlab.GetRawFileOptions{Ref: gitlab.String("master")})
		if err != nil {
			logger.Errorf("[%s] Downloading raw blob content to local...FAILED: %v", transID, res.Status, err)
			return "FAILED"
		}
		logger.Debugf("[%s] Downloading raw blob content to local...OK", transID, parentPath)

		writErr := os.WriteFile(rootPath + "/" + blobList[blob] + ".tmp", blobRawContent, 0644)
		if writErr != nil {
			logger.Errorf("[%s] Writing raw blob content to temporary file...%s: %v", transID, res.Status, writErr)
			return "FAILED"
		}
		logger.Debugf("[%s] Writing raw blob content to temporary file...%s", transID, parentPath)

		// Replace old tag with new tag
		replaceState := utils.ReplaceImageTag(transID, rootPath + "/" + blobList[blob], oldTag, newTag)
		if replaceState {
			logger.Debugf("[%s] Replacing old tag [%s] with new tag [%s]...OK", transID, oldTag, newTag)
		}

		if environment == "prod" {
			status := commitChange(projectID, transID, imageName, branchName, oldTag, newTag, blobList, blob, rootPath + "/" + blobList[blob])
			saveState(transID, fullImage, cluster, blobList[blob], status)
		} else {
			status := commitChange(projectID, transID, imageName, "master", oldTag, newTag, blobList, blob, rootPath + "/" + blobList[blob])
			saveState(transID, fullImage, cluster, blobList[blob], status)
		}
	}

	// Commit a new change
	if environment == "prod" {
		created := createMergeRequest(projectID, transID, imageName, branchName, newTag)
		if created { return "SUCCESSFUL" }
	}

	return "SUCCESSFUL"
}

func commitChange(projectID interface{}, transID string, imageName string, branchName string, oldTag string, newTag string, blobList []string, blob int, filePath string) string {
	logger := utils.ConfigZap()
	client := createNewGitlabClient()
	content, err := os.ReadFile(blobList[blob])
	if err != nil {
		logger.Errorf("[%s] Reading temporary file %s...FAILED: %v", transID, blob, err)
	}

	_, res, err := client.Commits.CreateCommit(projectID, &gitlab.CreateCommitOptions{
		Branch: gitlab.String(branchName),
		CommitMessage: gitlab.String("Change tag for " + imageName + " from oldtag " + oldTag + " to newtag " + newTag),
		Actions: []*gitlab.CommitActionOptions{
			{
				Action: gitlab.FileAction(gitlab.FileUpdate),
				FilePath: gitlab.String(blobList[blob]),
				Content: gitlab.String(string(content)),
			},
		},
	})
	if err != nil {
		logger.Errorf("[%s][%d] Committing new change to branch [%s]...%s: %v", transID, blob, branchName, res.Status, err)
		return "FAILED"
	}
	logger.Infof("[%s][%d] Committing new change to branch [%s]...%s", transID, blob, branchName, res.Status)
	return "SUCCESSFUL"
}

func createMergeRequest(projectID interface{}, transID string, imageName string, branchName string, newTag string) bool {
	logger := utils.ConfigZap()
	client := createNewGitlabClient()

	if branchName != "master" {
		_, res, err := client.MergeRequests.CreateMergeRequest(projectID,&gitlab.CreateMergeRequestOptions{
			Title: gitlab.String("Vnpaybot has released " + imageName + ":" + newTag),
			SourceBranch: gitlab.String(branchName),
			TargetBranch: gitlab.String("master"),
			AssigneeID: gitlab.Int(getUserID(projectID, transID)),
		})
		if err != nil {
			logger.Errorf("[%s] Creating a merge request for new branch %s...%s: %v", transID, branchName, res.Status, err)
			return false
		}
		logger.Infof("[%s] Creating a merge request for new branch %s...%s", transID, branchName, res.Status)
		return true
	}

	return false
}

func listBranch(projectID interface{}, oldBranchName string, transID string) bool {
	logger := utils.ConfigZap()
	client := createNewGitlabClient()

	_, _, err := client.Branches.ListBranches(projectID, &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page: 1,
		},
		Search: gitlab.String(oldBranchName),
	})
	if err != nil {
		logger.Warnf("[%s] Checking if old branch %s is existing...NONEXISTENT", transID, oldBranchName)
		return false
	}

	return true
}

func createNewBranch(projectID interface{}, transID string, oldBranchName string, branchName string) (error,*gitlab.Response) {
	client := createNewGitlabClient()

	branchIsExists := listBranch(projectID, oldBranchName, transID)
	if branchIsExists {
		deleteOldBranch(projectID, oldBranchName, transID)
	}

	_, res, err := client.Branches.CreateBranch(projectID,
		&gitlab.CreateBranchOptions{
			Branch: gitlab.String(branchName),
			Ref: gitlab.String("master"),
		})

	if res.StatusCode == 201 {
		return nil,res
	}
	return err,res
}

func deleteOldBranch(projectID interface{}, branchName string, transID string) bool {
	logger := utils.ConfigZap()
	client := createNewGitlabClient()

	res, err := client.Branches.DeleteBranch(projectID, branchName)
	if err != nil {
		logger.Warnf("[%s] Deleting old branch %s...%s: %v", transID, branchName, res.Status, err)
		return false
	}
	logger.Debugf("[%s] Deleting old branch %s...%s", transID, branchName, res.Status)
	return true
}

func createNewGitlabClient() *gitlab.Client {
	logger := utils.ConfigZap()
	gitlabURL := os.Getenv("GITLABURL")
	gitlabToken := os.Getenv("GITLABTOKEN")

	git, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabURL + "/api/v4"))
	if err != nil {
		logger.Errorf("Creating a new Gitlab client...failed: %v", err)
	} else {
		logger.Debug("Creating a new Gitlab client...OK")
	}

	return git
}

func getBlobContent(projectID interface{}, blobName string, transID string) string {
	logger := utils.ConfigZap()
	client := createNewGitlabClient()

	blob, res, err := client.RepositoryFiles.GetFile(
		projectID,
		blobName,
		&gitlab.GetFileOptions{
			Ref: gitlab.String("master"),
		},
	)
	if err != nil {
		logger.Errorf("[%s] Getting blob content %s...%s: %v", transID, blobName, res.Status, err)
	} else {
		logger.Debugf("[%s] Getting blob content %s...%s", transID, blobName, res.Status)
	}

	return blob.Content
}

func getOldTag(projectID interface{}, transID string, blobList []string, imageName string) (string, string) {
	var oldTagList []string
	logger := utils.ConfigZap()
	re := regexp.MustCompile(`(?:tag?:)\s+((t|d)-[a-z0-9]{8})|(m-(\d+\.)+\d+-[a-z0-9]{8})`)
	for blob := 0; blob < len(blobList); blob++ {
		blobName := blobList[blob]
		blobContent := getBlobContent(projectID, blobName, transID)

		byteBlobContent, err := b64.StdEncoding.DecodeString(blobContent)
		if err != nil {
			logger.Errorf("[%s] Base64 decoding file [%s]...failed: %v", transID, blobName, err)
		} else {
			logger.Debugf("[%s] Base64 decoding file [%s]...OK", transID, blobName)
		}

		// trimedBlobContent := strings.TrimSpace(string(byteBlobContent))
		regexTag := strings.ReplaceAll(re.FindString(string(byteBlobContent)), "tag:", "")
		trimedRegexTag := strings.TrimSpace(regexTag)

		if !utils.Contains(oldTagList, trimedRegexTag) {
			oldTagList = append(oldTagList, trimedRegexTag)
		}
	}

	if len(oldTagList) > 1 {
		return "", "More than one old tag found"
	}

	if len(oldTagList) == 0 {
		return "", "Old tag not found"
	}

	return oldTagList[0], ""
}

func getUserID(projectID interface{}, transID string) int {
	logger := utils.ConfigZap()
	client := createNewGitlabClient()

	user, res, err := client.Users.CurrentUser()
	if err != nil {
		logger.Errorf("[%s] Getting current user ID...%s: %v", transID, res.Status, err)
	}

	return user.ID
}

func locateBlob(projectID interface{}, transID string, imageName string) []string {
	var blobList []string
	logger := utils.ConfigZap()
	client := createNewGitlabClient()

	blobs, res, err := client.Search.BlobsByProject(
		projectID,
		imageName,
		&gitlab.SearchOptions{
			Ref: gitlab.String("master"),
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page: 1,
			},
		})
	if err != nil {
		logger.Errorf("[%s] Locating files from project ID %d...%s: %v", transID, projectID, res.Status, err)
	} else {
		logger.Infof("[%s] Locating files from project ID %d...%s", transID, projectID, res.Status)
	}

	for i := 0; i < len(blobs); i++ {
		blobName := blobs[i].Filename
		blobList = append(blobList, blobName)
	}

	return blobList
}