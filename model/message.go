package model

type Message struct {
	Cluster  string
	Image	string
}

type MessageStatus struct {
	Image		string	`bson:"image,omitempty"`
	Cluster		string	`bson:"cluster,omitempty"`
	BlobName 	string	`bson:"blob,omitempty"`
	Status		string	`bson:"status,omitempty"`
}