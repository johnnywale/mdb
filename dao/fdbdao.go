package dao

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
)

type FdbDao struct {
	db *fdb.Database
}

func NewFdbDao() *FdbDao {
	s := &FdbDao{}
	s.init()
	return s
}

func (s *FdbDao) init() {
	fdb.MustAPIVersion(520)
	db := fdb.MustOpenDefault()
	s.db = &db;
	//	db, err := fdb.Open("/home/pdcc/Desktop/golang/src/foundationdb/dao/fdb.cluster", []byte("DB"))
	//	if err != nil {
	//		panic(err)
	//	}
	//	s.db = &db
	//	fmt.Printf("init db==============%+v \n", s.db)
}

func (s *FdbDao) Clean(path []string) {

	directory, err := directory.CreateOrOpen(s.db, path, []byte("layer"))
	if err != nil {
		panic(err)
	}
	_, err = s.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		tr.Clear(directory)
		return nil, nil
	})

}

func (s *FdbDao) Load(path []string) {
	directory, _ := directory.CreateOrOpen(s.db, path, []byte("layer"))
	s.db.ReadTransact(func(rtr fdb.ReadTransaction) (interface{}, error) {
		ri := rtr.GetRange(directory, fdb.RangeOptions{
			Limit:   9999999,
			Mode:    fdb.StreamingModeSerial,
			Reverse: false,
		}).Iterator()
		i := 0
		for ri.Advance() {
			kv := ri.MustGet()
			key, err := directory.Unpack(kv.Key)
			if key != nil {
				i++
			}
			if err != nil {
				panic(err)
				return nil, err
			}
		}
		fmt.Printf("readed %v \n", i)
		return nil, nil
	})

}

type ExecuteFunc func(fdb.Transaction, directory.DirectorySubspace)

func (s *FdbDao) Execute(path []string, e ExecuteFunc) error {
	directory, err := directory.CreateOrOpen(s.db, path, []byte("layer"))
	if err != nil {
		return err
	}
	_, err = s.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		e(tr, directory)
		return nil, err
	})

	if err != nil {
		panic(err)
	}

	return nil
}

func (s *FdbDao) Save(path []string, key interface{}, protobuf []byte) error {
	directory, err := directory.CreateOrOpen(s.db, path, []byte("layer"))
	if err != nil {
		return err
	}
	_, err = s.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		tr.Set(directory.Pack(tuple.Tuple{key}), protobuf)
		return nil, nil
	})
	return nil
}
