# fileuploader

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
        <li><a href="#usage">Usage</a></li>  
        <li><a href="#example">Example</a></li> 
      </ul>
    </li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

This service is used to upload file in server. This service can be used for uploading config files needed for prometheus, consul etc.

## getting-started

### Built With
 golang
 
### installation
 
 step 1. download from <a href=https://github.com/smutil/fileuploader/releases>releases</a>. 
 
 step 2. start the service as shwon below
 
 ```
 ./fileuploader -dest /tmp/configs
 ```
 
### usage

``` 
  ./fileuploader -h
  -dest string
        (required) default destination directory
  -port string
        overwrite default port (default "3000")
 ```
 
 ### example

 1. uploading file in default destination directory
 ```
 curl -F "data=@test.yml"  http://localhost:3000/upload
 
 ```

 2. uploading file in user defined destination directory
 ```
 curl -F "data=@test.yml" -F "destDir=/tmp/downloads"  http://localhost:3000/upload
 
 ```
