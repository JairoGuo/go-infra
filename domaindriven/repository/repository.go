package repository

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32
}

type String interface {
	~string
}

type Identifier interface {
	Int | String
}

type Repository[E any, ID Identifier] interface {
	Save(entity E)
	Remove(id ID)
	Update(entity E) E
	Find(entity E) E
}
