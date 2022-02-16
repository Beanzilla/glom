# glom

This is based on the Python module/package [glom](https://pypi.org/project/glom/).

## Installation

Should be `go get github.com/beanzilla/glom` (Should automatically get [structs](https://pkg.go.dev/github.com/fatih/structs))

## Basic Example

```go
/* For brevity I will be showing a basic structure in a Pythonic format (Looks simular to JSON or HJSON)

// See Example Data section below for how this structure would/could look in Go.

data = { // Assume data is map[string]interface{}
    "users": {
    	"Bob": {
        	"last-on": "some time ago",
        	"age": 42,
        	"sec-level": 99
    	},
    	"Phil": {
	          "last-on": "a few seconds ago",
    	      "age": 0,
        	  "sec-level": 10
      	}
    },
    "posts": [
        {
            "title": "Why Glom-Go Rocks!",
            "post-date": "a few weeks ago",
            "description": "Access nested structures with ease, especially mixed types like map, slice/array, and interface.",
            "line-count": 1,
            "likes": 0
        },
        {
            "title": "Example of Glom-Go, and 5 other neat tips",
            "post-date": "a few months ago",
            "description": "Example in example... See recursion.",
            "line-count": 1,
            "likes": 0
        }
    ]
}
*/

// Let's start with something simple, getting the last time Bob was on...
bob_last_on, err := glom.Glom(data, "users.Bob.laston")
// In Python it would be data["users"]["Bob"]["laston"] (But in Go, we can't do that, due to our base type of data... interface)

if err != nil {
    fmt.Println(err) // Oops, and error occured, Looks like the error would be something like...
    // Failed moving to 'laston' from path of 'users.Bob', options are 'last-on', 'age', 'sec-level' (3)
    // So as the error is trying to let us know, we miss typed last-on with laston.
} else {
    fmt.Printf("Bob was last on %v.", bob_last_on) // If you fixed it so the string passed to glom.Glom was "users.Bob.last-on"
    // You now get:
    // Bob was last on some time ago.
}

// Now let's access the likes of the second post...
second_post_likes, err := glom.Glom(data, "posts.1.likes") // Remember slices/arrays start at 0, so the index of 1 will give us the second.

if err != nil {
    fmt.Println(err)
} else {
    fmt.Printf("Post 2 got %v likes.", second_post_likes)
    // Post 2 got 0 likes.
}
```



## Example Data

> This part was separated due to it's length and complexity. (It is added to allow the above example to actually work, just copy this then copy the Basic Example, above)

```go
// Initalize a multi-layer structure (I will try my best to show the layers deep, where 1 is directly accessable)
data := make(map[string]interface{})

// 1 can be like `data["users"]`
// 2 (which can't be directly accessed without some interface casting) can be `data["users"]["Bob"]`
// 3 (which can't be directly accessed without some nested interface casting) can be `data["users"]["Bob"]["sec-level"]`

users := make(map[string]interface{}) // 1

bob := make(map[string]interface{}) // 2
bob["last-on"] = "some time ago" // 3
bob["age"] = 42 // 3
bob["sec-level"] = 99 // 3
users["Bob"] = bob // Add Bob to users

phil := make(map[string]interface{}) // 2
phil["last-on"] = "a few seconds ago" // 3
phil["age"] = 0 // 3
phil["sec-level"] = 10 // 3
users["Phil"] = phil // Add Phil to users

data["users"] = users // Add users to data

var posts []interface{} // 1

post1 := make(map[string]interface{}) // 2
post1["title"] = "Why Glom-Go Rocks!" // 3
post1["post-date"] = "a few weeks ago" // 3
post1["description"] = "Access nested structures with ease, especially mixed types like map, slice/array, and interface." // 3
post1["line-count"] = 1 // 3
post1["likes"] = 0 // 3
posts = append(posts, post1) // Add post1 to posts

post2 := make(map[string]interface{}) // 2
post2["title"] = "Example of Glom-Go, and 5 other neat tips" // 3
post2["post-date"] = "a few months ago" // 3
post2["description"] = "Example in example... See recursion." // 3
post2["line-count"] = 1 // 3
post2["likes"] = 0 // 3
posts = append(posts, post2) // Add post2 to posts

data["posts"] = posts // Add posts to data

// data is now ready for glom.Glom

```

## Basic Accessing

Python's glom made accessing nested structures a breeze, glom was built to attempt to do just that just for Go.

* Just use dot notation for accessing/walking your data. (I.E. `users.Bob.age` would access the 3rd level, and comes with Python's glom error messaging showing exactly where while it walked the data it got lost at/couldn't go)
* Special star for dot notation. (I.E. to get all the users you could `users` or even `users.*`)
* glom can support maps, slices/arrays, and structures (specifically it's fields that are public/exposed), making it rather extensive.

## Under the hood

glom uses reflect and [structs](https://pkg.go.dev/github.com/fatih/structs) to handle nesting thru various structures.

From `[]interface{}` to `map[string]interface{}` (Even `map[int]interface{}` works).

All while supporting custom structures...

```go
type User struct {
    Name string
    Last_on string
    Age int
}

type Post struct {
    Title string
    Description string
    Line_count int
    Author User
    Post_date string
}

type Blog struct {
    Site_name string
    Posts []Post
    Site_owner User
}
bob := User{"Bob", "some time ago", 42}
blog := Blog{Site_name: "Test Site", Site_Owner: bob}
blog.Posts = append(blog.Posts, Post{"Example of Structs with glom", "Yet another example of Glom-Go", 1, bob, "a few seconds ago"})

// Example of accessing that to get first post's author
post_owner, err := glom.Glom(blog, "Posts.0.Author.Name")
if err != nil {
    fmt.Println(err)
} else {
    fmt.Printf("%v wrote the first post!", post_owner)
    // Bob wrote the first post!
}
```

