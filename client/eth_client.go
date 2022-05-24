// Copyright 2022 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"

	"github.com/ubiq/go-ubiq/v7/ethclient"
)

type EthClient struct {
	*ethclient.Client
}

// NewEthClient connects a SDKClient to the given URL.
func NewEthClient(endpoint string) (*EthClient, error) {
	client, err := ethclient.Dial(endpoint)

	if err != nil {
		return nil, fmt.Errorf("%w: unable to dial node", err)
	}

	return &EthClient{client}, nil
}

// Close shuts down the RPC SDKClient connection.
func (ec *EthClient) Close() {}
