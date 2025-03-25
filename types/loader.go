package types

type LoaderMessage struct {
	Message    string
	DelayAfter int
	Icon       string
	Color      string

	Progress        bool
	ProgressProcess interface{}
	ProgressTotal   interface{}
	ProgressCurrent interface{}
	ProgressMessage string
	ProgressIcon    string
}

type LoaderOrchestrator struct {
	LoaderMessages []LoaderMessage
}

func (l *LoaderOrchestrator) AddMessage(m LoaderMessage) {
	l.LoaderMessages = append(l.LoaderMessages, m)
}
func (l *LoaderOrchestrator) AddMessages(m []LoaderMessage) {
	l.LoaderMessages = append(l.LoaderMessages, m...)
}
func (l *LoaderOrchestrator) ClearMessages()               { l.LoaderMessages = make([]LoaderMessage, 0) }
func (l *LoaderOrchestrator) GetMessages() []LoaderMessage { return l.LoaderMessages }
func (l *LoaderOrchestrator) GetMessagesCount() int        { return len(l.LoaderMessages) }
func (l *LoaderOrchestrator) GetLastMessage() LoaderMessage {
	if len(l.LoaderMessages) > 0 {
		return l.LoaderMessages[len(l.LoaderMessages)-1]
	}
	return LoaderMessage{}
}
func (l *LoaderOrchestrator) GetFirstMessage() LoaderMessage {
	if len(l.LoaderMessages) > 0 {
		return l.LoaderMessages[0]
	}
	return LoaderMessage{}
}

func NewLoaderOrchestrator() *LoaderOrchestrator {
	return &LoaderOrchestrator{}
}

type Loader interface {
	AddMessage(LoaderMessage)
	AddMessages([]LoaderMessage)
	ClearMessages()
	GetMessages() []LoaderMessage
	GetMessagesCount() int
	GetLastMessage() LoaderMessage
	GetFirstMessage() LoaderMessage
}
