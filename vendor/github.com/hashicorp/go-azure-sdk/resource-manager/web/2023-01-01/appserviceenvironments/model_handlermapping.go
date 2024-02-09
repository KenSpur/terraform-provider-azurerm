package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HandlerMapping struct {
	Arguments       *string `json:"arguments,omitempty"`
	Extension       *string `json:"extension,omitempty"`
	ScriptProcessor *string `json:"scriptProcessor,omitempty"`
}
