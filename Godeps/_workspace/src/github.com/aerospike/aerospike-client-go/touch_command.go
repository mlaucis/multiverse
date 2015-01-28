// Copyright 2013-2014 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	. "github.com/aerospike/aerospike-client-go/types"
)

// guarantee touchCommand implements command interface
var _ command = &touchCommand{}

type touchCommand struct {
	singleCommand

	policy *WritePolicy
}

func newTouchCommand(cluster *Cluster, policy *WritePolicy, key *Key) *touchCommand {
	newTouchCmd := &touchCommand{
		singleCommand: *newSingleCommand(cluster, key),
	}

	if policy == nil {
		newTouchCmd.policy = NewWritePolicy(0, 0)
	} else {
		newTouchCmd.policy = policy
	}

	return newTouchCmd
}

func (cmd *touchCommand) getPolicy(ifc command) Policy {
	return cmd.policy
}

func (cmd *touchCommand) writeBuffer(ifc command) error {
	return cmd.setTouch(cmd.policy, cmd.key)
}

func (cmd *touchCommand) parseResult(ifc command, conn *Connection) error {
	// Read header.
	if _, err := conn.Read(cmd.dataBuffer, int(_MSG_TOTAL_HEADER_SIZE)); err != nil {
		return err
	}

	resultCode := cmd.dataBuffer[13] & 0xFF

	if resultCode != 0 {
		return NewAerospikeError(ResultCode(resultCode))
	}
	if err := cmd.emptySocket(conn); err != nil {
		return err
	}
	return nil
}

func (cmd *touchCommand) Execute() error {
	return cmd.execute(cmd)
}
