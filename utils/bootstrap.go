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

package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ubiq/rosetta-gubiq-sdk/configuration"
	"github.com/ubiq/rosetta-gubiq-sdk/services"
	"github.com/ubiq/rosetta-gubiq-sdk/services/construction"

	AssetTypes "github.com/ubiq/rosetta-gubiq-sdk/types"

	"github.com/neilotoole/errgroup"
	"github.com/ubiq/rosetta-sdk-go/asserter"
	"github.com/ubiq/rosetta-sdk-go/server"
	RosettaTypes "github.com/ubiq/rosetta-sdk-go/types"
)

// BootStrap quickly starts the Rosetta server
// and begin to serve Rosetta RESTful requests
func BootStrap(
	cfg *configuration.Configuration,
	types *AssetTypes.Types,
	errors []*RosettaTypes.Error,
	client construction.Client,
) error {
	// The asserter automatically rejects incorrectly formatted
	// requests.
	asserter, err := asserter.NewServer(
		types.OperationTypes,
		AssetTypes.HistoricalBalanceSupported,
		[]*RosettaTypes.NetworkIdentifier{cfg.Network},
		types.CallMethods,
		AssetTypes.IncludeMempoolCoins,
		"",
	)
	if err != nil {
		return fmt.Errorf("%w: could not initialize server asserter", err)
	}
	router := services.NewBlockchainRouter(cfg, types, errors, client, asserter)

	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: corsRouter,
	}

	// Start required services
	ctx := context.Background()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Printf("server listening on port %d", cfg.Port)
		return server.ListenAndServe()
	})

	g.Go(func() error {
		// If we don't shutdown server in errgroup, it will
		// never stop because server.ListenAndServe doesn't
		// take any context.
		<-ctx.Done()

		return server.Shutdown(ctx)
	})

	err = g.Wait()

	return err
}
