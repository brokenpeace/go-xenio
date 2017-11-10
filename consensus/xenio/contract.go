// Copyright 2017 The go-xenio Authors
//
// This file is part of the go-xenio library.
//
// The go-xenio library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-xenio library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-xenio library. If not, see <http://www.gnu.org/licenses/>.
//
// This file provides generic functions to interact with contracts

package xenio

import (
	"io/ioutil"
	"strings"
	"strconv"
	"errors"
	"sync"

	"github.com/xenioplatform/go-xenio/accounts/abi/bind"
	"github.com/xenioplatform/go-xenio/accounts/keystore"
	"github.com/xenioplatform/go-xenio/common"
	"github.com/xenioplatform/go-xenio/ethclient"
	"github.com/xenioplatform/go-xenio/log"
)

// xenio contract related variables (users, servers and games).
var (
	currentIPCEndpoint    string // required for interaction with contracts
	currentKeyStoreDir    string // required for paid contract interactions
	deployedGamesContract common.Address
	deployedUsersContract common.Address

	currentTransactor transactor // authorized signer for paid transactions
)

var (
	// Contract Specific Errors
	errTransactorNotSet  = errors.New("transactor not set")
	//errTransactorExpired = errors.New("transactor expired")
)

type transactor struct {
	contractAuth           *bind.TransactOpts
	authorizedTransactions int

	lock   sync.RWMutex   // Protects the transactor fields
}

func (api *API) SetContractTransactor(fromAddress common.Address, pwd string, transactionsTL int) error {
	ks := keystore.NewKeyStore(currentKeyStoreDir, keystore.LightScryptN, keystore.LightScryptP)
	localAccounts := ks.Accounts()
	keyjson, err := ioutil.ReadFile(localAccounts[0].URL.Path)
	if &transactionsTL == nil || transactionsTL == 0 {
		transactionsTL = 1
	}
	currentTransactor.authorizedTransactions = transactionsTL
	currentTransactor.contractAuth, err = bind.NewTransactor(strings.NewReader(string(keyjson)), pwd)

	log.Info("Set auth for " + strconv.Itoa(transactionsTL) + " transactions")
	return err
}

func resetContractTransactor() {
	if currentTransactor.authorizedTransactions > 0 {
		currentTransactor.authorizedTransactions = currentTransactor.authorizedTransactions - 1
	}
	// transactor expired
	if currentTransactor.authorizedTransactions == 0 {
	currentTransactor = transactor{}
	}
}

func getContractBackend() (*ethclient.Client, error) {
	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial(currentIPCEndpoint)
	if err != nil {
		log.Error("Failed to connect to the Xenio client: " + err.Error())
	}
	return conn, err
}

func getFreeTxTransactor() *bind.CallOpts{
	var freeTxTransactor *bind.CallOpts
	return freeTxTransactor
}