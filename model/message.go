package model

type Message struct {
	Cluster  string
	Image	string
}

type MessageStatus struct {
	Image		string	`bson:"image,omitempty"`
	OldTag		string	`bson:"oldtag,omitempty"`
	NewTag		string	`bson:"newtag,omitempty"`
	Cluster		string	`bson:"cluster,omitempty"`
	BlobName 	string	`bson:"blob,omitempty"`
	Time 		string	`bson:"time,omitempty"`
	Status		string	`bson:"status,omitempty"`
}