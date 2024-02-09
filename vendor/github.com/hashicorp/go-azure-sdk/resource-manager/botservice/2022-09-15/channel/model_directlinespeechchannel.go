package channel

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Channel = DirectLineSpeechChannel{}

type DirectLineSpeechChannel struct {
	Properties *DirectLineSpeechChannelProperties `json:"properties,omitempty"`

	// Fields inherited from Channel
	Etag              *string `json:"etag,omitempty"`
	Location          *string `json:"location,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

var _ json.Marshaler = DirectLineSpeechChannel{}

func (s DirectLineSpeechChannel) MarshalJSON() ([]byte, error) {
	type wrapper DirectLineSpeechChannel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DirectLineSpeechChannel: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DirectLineSpeechChannel: %+v", err)
	}
	decoded["channelName"] = "DirectLineSpeechChannel"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DirectLineSpeechChannel: %+v", err)
	}

	return encoded, nil
}
