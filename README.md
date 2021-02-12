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

This service can be used to upload file. 

## getting-started

### Built With
 golang
 
### installation
 
 download from releases and copy into any server. 
 
 start the service as shwon below
 
 ```
 ./fileuploader -dest /tmp/configs
 ```
 
### usage

``` 
  ./fileuploader -h
  -dest string
        (required) destination directory
  -port string
        overwrite default port (default "3000")
 ```
 
 ### example
 ```
 curl -F "data=@test.yml"  http://localhost:3000/upload
 
  ```
