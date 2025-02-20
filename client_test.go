/*
 * MIT License
 *
 * Copyright (c) 2024 Arsene Tochemey Gandote
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package gokv

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/tochemey/gokv/internal/lib"
	"github.com/tochemey/gokv/test/data/testpb"
)

func TestClient(t *testing.T) {
	t.Run("With PutProto GetProto", func(t *testing.T) {
		ctx := context.Background()
		// start the NATS server
		srv := startNatsServer(t)
		// create a cluster node1
		node1, sd1 := startNode(t, srv.Addr().String())
		require.NotNil(t, node1)
		// create a cluster node2
		node2, sd2 := startNode(t, srv.Addr().String())
		require.NotNil(t, node2)
		key := "my-key"
		value := new(testpb.Hello)
		err := node2.Client().PutProto(ctx, key, value, NoExpiration)
		require.NoError(t, err)
		// wait for the key to be distributed in the cluster
		lib.Pause(time.Second)
		// let us retrieve the key from the other nodes
		exists, err := node1.Client().Exists(ctx, key)
		require.NoError(t, err)
		require.True(t, exists)
		actual := &testpb.Hello{}
		err = node1.Client().GetProto(ctx, key, actual)
		require.NoError(t, err)
		assert.True(t, proto.Equal(value, actual))
		lib.Pause(time.Second)

		assert.NoError(t, node1.Stop(ctx))
		assert.NoError(t, node2.Stop(ctx))
		assert.NoError(t, sd1.Close())
		assert.NoError(t, sd2.Close())
		srv.Shutdown()
	})
	t.Run("With PutString GetString", func(t *testing.T) {
		ctx := context.Background()
		// start the NATS server
		srv := startNatsServer(t)
		// create a cluster node1
		node1, sd1 := startNode(t, srv.Addr().String())
		require.NotNil(t, node1)
		// create a cluster node2
		node2, sd2 := startNode(t, srv.Addr().String())
		require.NotNil(t, node2)
		key := "my-key"
		value := "my-value"
		err := node2.Client().PutString(ctx, key, value, NoExpiration)
		require.NoError(t, err)
		// wait for the key to be distributed in the cluster
		lib.Pause(time.Second)
		// let us retrieve the key from the other nodes
		exists, err := node1.Client().Exists(ctx, key)
		require.NoError(t, err)
		require.True(t, exists)
		actual, err := node1.Client().GetString(ctx, key)
		require.NoError(t, err)
		require.NotEmpty(t, actual)
		require.Equal(t, value, actual)
		lib.Pause(time.Second)

		assert.NoError(t, node1.Stop(ctx))
		assert.NoError(t, node2.Stop(ctx))
		assert.NoError(t, sd1.Close())
		assert.NoError(t, sd2.Close())
		srv.Shutdown()
	})
	t.Run("With PutProto GetProto with expiration", func(t *testing.T) {
		ctx := context.Background()
		// start the NATS server
		srv := startNatsServer(t)
		// create a cluster node1
		node1, sd1 := startNode(t, srv.Addr().String())
		require.NotNil(t, node1)
		// create a cluster node2
		node2, sd2 := startNode(t, srv.Addr().String())
		require.NotNil(t, node2)
		expiration := 100 * time.Millisecond
		key := "my-key"
		value := &testpb.Hello{Name: key}
		err := node2.Client().PutProto(ctx, key, value, expiration)
		require.NoError(t, err)
		// wait for the key to be distributed in the cluster
		lib.Pause(time.Second)
		// let us retrieve the key from the other nodes
		exists, err := node1.Client().Exists(ctx, key)
		require.NoError(t, err)
		require.False(t, exists)
		actual := &testpb.Hello{}
		err = node1.Client().GetProto(ctx, key, actual)
		require.Error(t, err)
		assert.EqualError(t, err, ErrKeyNotFound.Error())
		lib.Pause(time.Second)
		assert.NoError(t, node1.Stop(ctx))
		assert.NoError(t, node2.Stop(ctx))
		assert.NoError(t, sd1.Close())
		assert.NoError(t, sd2.Close())
		srv.Shutdown()
	})
	t.Run("With PutString GetString with expiration", func(t *testing.T) {
		ctx := context.Background()
		// start the NATS server
		srv := startNatsServer(t)
		// create a cluster node1
		node1, sd1 := startNode(t, srv.Addr().String())
		require.NotNil(t, node1)
		// create a cluster node2
		node2, sd2 := startNode(t, srv.Addr().String())
		require.NotNil(t, node2)
		key := "my-key"
		value := "my-value"
		expiration := 100 * time.Millisecond
		err := node2.Client().PutString(ctx, key, value, expiration)
		require.NoError(t, err)
		// wait for the key to be distributed in the cluster
		lib.Pause(time.Second)
		// let us retrieve the key from the other nodes
		exists, err := node1.Client().Exists(ctx, key)
		require.NoError(t, err)
		require.False(t, exists)
		actual, err := node1.Client().GetString(ctx, key)
		require.Error(t, err)
		require.Empty(t, actual)
		assert.EqualError(t, err, ErrKeyNotFound.Error())
		lib.Pause(time.Second)
		assert.NoError(t, node1.Stop(ctx))
		assert.NoError(t, node2.Stop(ctx))
		assert.NoError(t, sd1.Close())
		assert.NoError(t, sd2.Close())
		srv.Shutdown()
	})
	t.Run("With Put Get with expiration", func(t *testing.T) {
		ctx := context.Background()
		// start the NATS server
		srv := startNatsServer(t)
		// create a cluster node1
		node1, sd1 := startNode(t, srv.Addr().String())
		require.NotNil(t, node1)

		// create a cluster node2
		node2, sd2 := startNode(t, srv.Addr().String())
		require.NotNil(t, node2)

		expiration := 100 * time.Millisecond
		key := "my-key"
		value := &testpb.Hello{Name: key}
		bytea, err := proto.Marshal(value)
		require.NoError(t, err)
		entry := &Entry{
			Key:   key,
			Value: bytea,
		}

		err = node2.Client().Put(ctx, entry, expiration)
		require.NoError(t, err)

		// wait for the key to be distributed in the cluster
		lib.Pause(time.Second)

		// let us retrieve the key from the other nodes
		exists, err := node1.Client().Exists(ctx, key)
		require.NoError(t, err)
		require.False(t, exists)

		actual, err := node1.Client().Get(ctx, key)
		require.Error(t, err)
		require.Nil(t, actual)
		assert.EqualError(t, err, ErrKeyNotFound.Error())

		lib.Pause(time.Second)
		t.Cleanup(func() {
			assert.NoError(t, node1.Stop(ctx))
			assert.NoError(t, node2.Stop(ctx))
			assert.NoError(t, sd1.Close())
			assert.NoError(t, sd2.Close())
			srv.Shutdown()
		})
	})
}

type testCodec struct{}

func (c *testCodec) Encode(t interface{}) ([]byte, error) {
	return json.Marshal(t)
}
func (c *testCodec) Decode(bytea []byte) (interface{}, error) {
	var t interface{}
	err := json.Unmarshal(bytea, &t)
	return t, err
}
