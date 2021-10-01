package store

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"
	"passwordSafe/internal"
)

const KEYSIZE = 4
const VALSIZE = 4

type TxtStore struct {
	credsFile *os.File
}

func NewTxtStore(credsFile *os.File) internal.Store {
	return TxtStore{credsFile: credsFile}
}

func (s TxtStore) Save(c internal.Credentials) error {
	//4 - max bytes for key length
	var keyLen []byte = make([]byte, KEYSIZE)
	var valLen []byte = make([]byte, VALSIZE)
	for k, v := range c {
		if len(k) <= math.MaxUint32 {
			binary.LittleEndian.PutUint32(keyLen, uint32(len(k)))
		} else {
			return errors.New("key is too large")
		}

		if len(v) <= math.MaxUint32 {
			binary.LittleEndian.PutUint32(valLen, uint32(len(v)))
		} else {
			return errors.New("val is too large")
		}

		n, err := s.credsFile.Write(append(append(append(keyLen, valLen...), []byte(k)...), v...))
		if err != nil {
			return err
		}

		if n != len(valLen) + len(keyLen) + len([]byte(k)) + len(v) && n != 0 {
			i, err := s.credsFile.Stat()
			if err != nil {
				return err
			}

			err = os.Truncate(s.credsFile.Name(), i.Size()-int64(n))
			if err != nil {
				return err
			}
			return errors.New("notation has failed to wrote down, rollback")
		}
	}
	return nil
}

func (s TxtStore) Get() (internal.Credentials, error) {
	_, err := s.credsFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var creds = make(internal.Credentials)
	var keySize = make([]byte, KEYSIZE)
	var valSize = make([]byte, VALSIZE)
	for {
		//READ KEY SIZE
		n, err := s.credsFile.Read(keySize)
		if err != nil {
			if err == io.EOF {
				break
				//return nil, errors.New("EOF error while reading keySize think what to do ||| TODO")
			}
			return nil, err
		}

		if n != KEYSIZE {
			return nil, errors.New("cant read key size properly")
		}

		keySizeUint := binary.LittleEndian.Uint32(keySize)
		key := make([]byte, keySizeUint)

		//READ VAL SIZE
		n, err = s.credsFile.Read(valSize)
		if err != nil {
			if err == io.EOF {
				break
				//return nil, errors.New("EOF error while reading keySize think what to do ||| TODO")
			}
			return nil, err
		}

		if n != VALSIZE {
			return nil, errors.New("cant read key size properly")
		}

		valSizeUint := binary.LittleEndian.Uint32(valSize)
		val := make([]byte, valSizeUint)

		//READ KEY
		n, err = s.credsFile.Read(key)
		if err != nil {
			if err == io.EOF {
				break
				//return nil, errors.New("EOF error while reading key think what to do ||| TODO")
			}
			return nil, err
		}

		if n != len(key) {
			return nil, errors.New("cant read key properly")
		}

		//READ VAL
		n, err = s.credsFile.Read(val)
		if err != nil {
			if err == io.EOF {
				break
				//return nil, errors.New("EOF error while reading val think what to do ||| TODO")
			}
			return nil, err
		}

		if n != len(val) {
			return nil, errors.New("cant read val properly")
		}

		//SAVE TO MAP
		creds[string(key)] = val
	}

	return creds, nil
}