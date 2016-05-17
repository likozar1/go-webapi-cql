package main

import(
    "time"
    "fmt"
    "github.com/gocql/gocql"
    "log"
    "strconv"
 )

const(
    HOST          = "127.0.0.1"
    KEYSPACE      = "fileservice"
    TABLE_NAME    = "files"
    TIMEOFFSET    = "-0700"
)

var data Data
var currentId int

// Give us some seed data

func RepoFindFile(id int) *File {
    ids:=strconv.Itoa(id)
    file:=dbGet(ids)
    data.append(file)
    // return empty File if not found
    return file
}

func RepoFindAllFiles(){
    dbGetAll(&data)
}

func RepoCreateFile(t *File) *File {
    currentId:=dbCount()
    currentId += 1
    t.Id = currentId
    t.Due = time.Now()
    data.append(t)
    dbInsert(t)
    return t
}

func RepoDestroyTodo(id int) error {
    for i, t := range data.files {
        if t.Id == id {
            data.files = append(data.files[:i], data.files[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}

func dbGet(id string)*File{
    // connect to the cluster
    cluster := gocql.NewCluster(HOST)
    cluster.Keyspace = KEYSPACE
    session, _ := cluster.CreateSession()

    defer session.Close()
    return getFile(session,id)
}

//Fetch data from db
func getFile(session *gocql.Session,idd string) *File{
    file:=&File{}
    var (
        id int
        completed bool
        name string
        due time.Time
    )   
    query := session.Query("SELECT id,name,completed,due FROM "+TABLE_NAME+" where id="+idd+";").Iter()
    for query.Scan(&id,&name,&completed,&due) {
        file.Id=id
        file.Name=name
        file.Completed=completed
        file.Due=due
    }
    return file
}

func dbInsert(file *File){
    // connect to the cluster
    cluster := gocql.NewCluster(HOST)
    cluster.Keyspace = KEYSPACE
    session, _ := cluster.CreateSession()
    
    insertFile(session,file)

    defer session.Close()
}

//Fetch data from db
func insertFile(session *gocql.Session,file *File){  
    id:=strconv.Itoa(file.Id)
    bo:="false"
    if file.Completed==true{
        bo="true"
    }
    if err := session.Query("INSERT INTO "+TABLE_NAME+" (id,completed,due,name) VALUES ("+id+", "+bo+",'"+timeFormat(file.Due)+"','"+file.Name+"')").Exec(); err != nil {
        log.Fatal(err)
    }
}

func dbCount() int{
    // connect to the cluster
    cluster := gocql.NewCluster(HOST)
    cluster.Keyspace = KEYSPACE
    session, _ := cluster.CreateSession()
    
    

    defer session.Close()

    return CountFile(session)
}

//Fetch data from db
func CountFile(session *gocql.Session) int{  
    var (
        id int
    )   
    query := session.Query("SELECT COUNT(*) FROM "+TABLE_NAME+";").Iter()
    for query.Scan(&id) {

    }
    return id
}

func dbGetAll(data *Data){
    // connect to the cluster
    cluster := gocql.NewCluster(HOST)
    cluster.Keyspace = KEYSPACE
    session, _ := cluster.CreateSession()
    getAllFile(session,data)
    defer session.Close()
}

//Fetch data from db
func getAllFile(session *gocql.Session,data *Data){
    
    var (
        id int
        completed bool
        name string
        due time.Time
    )   
    query := session.Query("SELECT id,name,completed,due FROM "+TABLE_NAME+";").Iter()
    for query.Scan(&id,&name,&completed,&due) {
        file:=&File{}
        file.Id=id
        file.Name=name
        file.Completed=completed
        file.Due=due
        data.append(file)


    }
}

func timeFormat(ts time.Time)string{
    return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d%s", ts.Year(), 
        ts.Month(), ts.Day(), ts.Hour(), 
        ts.Minute(), ts.Second(),ts.Format(TIMEOFFSET))
}