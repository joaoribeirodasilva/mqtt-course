package mqtt_client

type Topic struct {
	Name       string `json:"name"`
	Qos        byte   `json:"qos"`
	Retain     bool   `json:"retain"`
	Subscribed bool   `json:"subscribed"`
	Disabled   bool   `json:"subscribed"`
}

type TopicsMap struct {
	keys   []string
	pos    int
	topics map[string]*Topic
}

type Topics struct {
	pos int
	// TODO: map with topic names ?
	publisherTopics  TopicsMap
	subscriberTopics TopicsMap
	persistPath      string
}

func NewTopics(persistPath string) *Topics {

	t := &Topics{
		pos:              0,
		publisherTopics:  make([]*Topic, 0),
		subscriberTopics: make([]*Topic, 0),
		persistPath:      persistPath,
	}

	return t
}

func (t *Topics) AppendPublisher(topic *Topic) {

	t.publisherTopics = append(t.publisherTopics, topic)
}

func (t *Topics) GetNext() *Topic {

	if len(t.topics) < t.pos {
		return nil
	}

	t.pos++
	return t.topics[t.pos-1]
}

func (t *Topics) Load() error {

	// load file
	// parse json
	// unmarshal json

	return nil
}

func (t *Topics) Save() error {

	// check directory
	// marshal json
	// save file

	return nil
}

func NewTopic(name string, qos byte, retain bool, subscribed bool) *Topic {

	t := &Topic{
		Name:       name,
		Qos:        qos,
		Retain:     retain,
		Subscribed: subscribed,
	}

	return t
}
