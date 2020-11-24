Project Skeleton
================

Project skeleton based on our latest standards.
Provides a basic golang based rest service with redis, elastic, etc.

## Usage
### Creating a new project
1. `cd generate`
2. `cp meta.example.yml meta.yml`
3. modify meta.yml as needed
4. `go run generate.go`
5. copy content of output directory
6. run `go mod tidy` in the new project

### Updating an existing project
1. `cd generate`
2. copy meta.yml from the project to the generate directory
3. modify meta.yml as needed
4. `go run generate.go`
5. compare output directory with the existing project, and apply necessary changes
