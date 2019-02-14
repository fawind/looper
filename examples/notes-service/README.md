# Example Notes Micro-Services

Example micro-services used for testing.


```
                 -------------------
                 | Note-ID Service |
                 -------------------
                /
----------------
| Note Service |
----------------
                \
                ------------------------
                | Note Storage Service |
                ------------------------

```

## Note Service API

```
Note {
  id: int,
  content: string,
}

NewNote {
  content: string,
}
```

```
GET /notes -> Note[]
```

```
POST /notes -> NewNote -> Note
```

## Note-ID Service API

```
GET /new -> int
```

## Note Storage API

```
GET /notes -> Note[]

POST /notes -> Note -> void
```
