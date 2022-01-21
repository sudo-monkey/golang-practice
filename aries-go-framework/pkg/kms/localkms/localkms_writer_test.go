/*
 Copyright SecureKey Technologies Inc. All Rights Reserved.

 SPDX-License-Identifier: Apache-2.0
*/

package localkms

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/google/tink/go/subtle/random"
	"github.com/stretchr/testify/require"

	"github.com/hyperledger/aries-framework-go/pkg/kms"
	mockstorage "github.com/hyperledger/aries-framework-go/pkg/mock/storage"
)

func TestLocalKMSWriter(t *testing.T) {
	t.Run("success case - create a valid storeWriter and store 20 non empty random keys", func(t *testing.T) {
		storeMap := map[string]mockstorage.DBEntry{}
		mockStore := &mockstorage.MockStore{Store: storeMap}

		for i := 0; i < 256; i++ {
			l := newWriter(mockStore)
			require.NotEmpty(t, l)
			someKey := random.GetRandomBytes(uint32(32))
			n, err := l.Write(someKey)
			require.NoError(t, err)
			require.Equal(t, len(someKey), n)
			require.Equal(t, maxKeyIDLen, len(l.KeysetID), "for key creation iteration %d", i)
			retrievedDBEntry, ok := storeMap[l.KeysetID]
			require.True(t, ok)
			require.Equal(t, retrievedDBEntry.Value, someKey)
		}

		require.Equal(t, 256, len(storeMap))
	})

	t.Run("error case - create a storeWriter using a bad storeWriter to storage", func(t *testing.T) {
		storeMap := map[string]mockstorage.DBEntry{}
		putError := fmt.Errorf("failed to put data")
		mockStore := &mockstorage.MockStore{
			Store:  storeMap,
			ErrPut: putError,
		}

		l := newWriter(mockStore)
		require.NotEmpty(t, l)
		someKey := []byte("someKeyData")
		n, err := l.Write(someKey)
		require.EqualError(t, err, putError.Error())
		require.Equal(t, 0, n)
	})

	t.Run("error case - create a storeWriter using a bad storeReader from storage", func(t *testing.T) {
		storeMap := map[string]mockstorage.DBEntry{}
		getError := fmt.Errorf("failed to get data")
		mockStore := &mockstorage.MockStore{
			Store:  storeMap,
			ErrGet: getError,
		}

		l := newWriter(mockStore)
		require.NotEmpty(t, l)
		someKey := []byte("someKeyData")
		n, err := l.Write(someKey)
		require.EqualError(t, err, getError.Error())
		require.Equal(t, 0, n)
	})

	t.Run("error case - create a storeWriter with a keysetID using a bad storeReader from storage",
		func(t *testing.T) {
			storeMap := map[string]mockstorage.DBEntry{}
			getError := fmt.Errorf("failed to get data")
			mockStore := &mockstorage.MockStore{
				Store:  storeMap,
				ErrGet: getError,
			}

			l := newWriter(mockStore, kms.WithKeyID(base64.RawURLEncoding.EncodeToString(random.GetRandomBytes(
				uint32(base64.RawURLEncoding.DecodedLen(maxKeyIDLen))))))

			require.NotEmpty(t, l)
			someKey := []byte("someKeyData")
			n, err := l.Write(someKey)
			require.EqualError(t, err, "got error while verifying requested ID: "+getError.Error())
			require.Equal(t, 0, n)
		})

	t.Run("error case - import duplicate keysetID", func(t *testing.T) {
		storeMap := map[string]mockstorage.DBEntry{}
		mockStore := &mockstorage.MockStore{Store: storeMap}

		l := newWriter(mockStore)
		require.NotEmpty(t, l)
		someKey := random.GetRandomBytes(uint32(32))
		n, err := l.Write(someKey)
		require.NoError(t, err)
		require.Equal(t, len(someKey), n)
		require.Equal(t, maxKeyIDLen, len(l.KeysetID))
		retrievedDBEntry, ok := storeMap[l.KeysetID]
		require.True(t, ok)
		require.Equal(t, retrievedDBEntry.Value, someKey)

		require.Equal(t, 1, len(storeMap))

		// create s second writer with keysetID created above
		l2 := newWriter(mockStore, kms.WithKeyID(l.KeysetID))

		_, err = l2.Write(someKey)
		require.EqualError(t, err, fmt.Sprintf("requested ID '%s' already exists, cannot write keyset",
			l.KeysetID))
	})
}
