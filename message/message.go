package message

type Message struct {
	From      string
	Previous  string
	Seq       string
	Timestamp string

	Content   string

	Hash      string
	Signature string
}
