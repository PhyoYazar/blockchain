// Package disk implements the ability to read and write blocks to disk
// writing each block to a separate block numbered file.
package disk

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/PhyoYazar/blockchain/foundation/blockchain/database"
)

// Disk represents the serialization implementation for reading and storing
// blocks in their own separate files on disk. This implements the database.Storage
// interface.
type Disk struct {
	dbPath string
}

// New constructs an Disk value for use.
func New(dbPath string) (*Disk, error) {
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, err
	}

	return &Disk{dbPath: dbPath}, nil
}

// Close in this implementation has nothing to do since a new file is
// written to disk for each now block and then immediately closed.
func (d *Disk) Close() error {
	return nil
}

// Write takes the specified database blocks and stores it on disk in a
// file labeled with the block number.
func (d *Disk) Write(blockData database.BlockData) error {

	// Marshal the block for writing to disk in a more human readable format.
	data, err := json.MarshalIndent(blockData, "", "  ")
	if err != nil {
		return err
	}

	// Create a new file for this block and name it based on the block number.
	f, err := os.OpenFile(d.getPath(blockData.Header.Number), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the new block to disk.
	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

// getPath forms the path to the specified block.
func (d *Disk) getPath(blockNum uint64) string {
	name := strconv.FormatUint(blockNum, 10)
	return path.Join(d.dbPath, fmt.Sprintf("%s.json", name))
}
