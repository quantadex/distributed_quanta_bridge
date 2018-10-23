package kv_store

/**
 * KVStore
 *
 * This is a simple key-value store that is thread-safe.
 * Tables (namespaces) are used in the store.
 *
 * The keys within a table are UNIQUE. A value cannot be nil.
 * Nil is returned when a value does not exist.
 *
 * Proper use for updating a key is to first do a read and get the old value.
 * When setting the new value you pass the old value.
 * The operation ONLY succeeds if the value has not changed.
 */
type KVStore interface {
	/*
	 * Connect
	 *
	 * Inputs:
	 *  name string - name of key value store
	 *
	 * Outputs:
	 *  - on success nil
	 *  - on failure error
	 *
	 * Establish a connection to specified KV store
	 */
	Connect(name string) error

	/*
	 * CreateTable
	 *
	 * Inputs:
	 *  tableName string - name of table to create
	 *
	 * Outputs:
	 *  - on success nil
	 *  - on failure error
	 *
	 * Creates a new table (name space) within KV store
	 *
	 */
	CreateTable(tableName string) error

	/*
	 * GetValue
	 *
	 * Inputs:
	 *  tableName string - name of table
	 *  key string - key to retrieve
	 *
	 * Outputs:
	 *  value string - value retrieved if successful. Nil if key not found or error
	 *  err error - Nil if successful or key not found. Otherwise propogate error.
	 *
	 * Retrieves the value associated with the given key. If access error encountered error is
	 * propogated. If the key is not found both value and error are nil
	 */
	GetValue(tableName string, key string) (value *string, err error)

	RemoveKey(tableName string, key string) error

	Put(tableName string, key []byte, value []byte) error

	/*
	 * GetAllValues
	 *
	 * Return all KV in the table as a map
	 */
	GetAllValues(tableName string) (map[string]string, error)

	/*
	 * SetValue
	 *
	 * Inputs:
	 *  tableName string - name of table
	 *  key string - key to update
	 *  oldValue - the previous value of the key as belived by caller
	 *  newValue - the new value to set for key
	 *
	 * Outputs:
	 *  - nil if successful
	 *  - error propopgated if error encountered
	 *
	 * Sets the key to the newValue provided that the caller correctly knew the previous state of the
	 * key as passed via "oldValue". If caller believes the key does not already exist use nil as oldValue.
	 */
	SetValue(tableName string, key string, oldValue string, newValue string) error
}

func NewKVStore() (KVStore, error) {
	return &BoltStore{}, nil
}
