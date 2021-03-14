# Henkernieuws

![Badge showing that the project is MIT licensed](https://img.shields.io/github/license/theonejonahgold/progressive-web-apps-2021?label=lijsens&style=flat-square) ![Badge showing amount of open issues](https://img.shields.io/github/issues/theonejonahgold/progressive-web-apps-2021?label=issjoes&style=flat-square)

[Live link](https://aqueous-beach-16784.herokuapp.com/)

Henkernieuws hes de best nieuws voor riel henkers lijk joe en mie.

Ritten in go en javascript, lijk ee pro.

## Table of contents

- [Getting started](#getting-started)
  - [Installing the project](#installing-the-project)
  - [Available commands](#available-commands)
- [Features](#features)


## Getting started

This program is built with Go@^1.16. This means that you have to install one of the Go 1.16 versions via the [golang website](https://golang.org). Also make sure you have [NodeJS](https://nodejs.org) installed with [Yarn v1](http://classic.yarnpkg.com) or (non-preferably) the [NPM package mananger](https://npmjs.com).

### Installing the project

```sh
$ git clone https://github.com/theonejonahgold/progressive-web-apps-2021.git pwa
$ cd pwa
$ go get && yarn # or "npm i"
$ go build -o pwa main.go
```

With this sequence of commands, you create a binary called `pwa` in the project root. While working on Henkernieuws, you use this binary to interact with your dev environment. **Do not use NPM scripts. They do not work**. If you want the NPM scripts to run, run `go install` first.

### Available commands

```sh
$ ./pwa start # Serves the contents of the dist folder, only useful if you've already prerendered the site.
$ ./pwa build # Downloads all henkernieuws data, and prerenders the entire site, also running snowpack to build everything.
$ ./pwa dev   # Runs a Fiber (express-like) server dynamically serving handlebars templates.
```

## Features

- [x] Static file serving
- [x] Prerendering as build step
- [x] Separate dynamic dev server
- [x] Best HackerNews clone to ever exist
- [ ] Client side routing
- [ ] Caching with service workers
- [ ] Custom JS API functions
- [ ] Periodic rerendering on production
- [ ] Making Node.JS unnecessary as a requirement

<!-- ### Week 1 - Server Side Rendering 📡

Goal: Render web pages server side

[Exercises](https://github.com/cmda-minor-web/progressive-web-apps-2021/blob/master/course/week-1.md)    
[Server Side Rendering - slides Declan Rek](https://github.com/cmda-minor-web/progressive-web-apps-1920/blob/master/course/cmd-2021-server-side-rendering.pdf)  


### Week 2 - Progressive Web App 🚀

Goals: Convert application to a Progressive Web App

[Exercises](https://github.com/cmda-minor-web/progressive-web-apps-2021/blob/master/course/week-2.md)  
[Progressive Web Apps - slides Declan Rek](https://github.com/cmda-minor-web/progressive-web-apps-1920/blob/master/course/cmd-2020-progressive-web-apps.pdf)


### Week 3 - Critical Rendering Path 📉 

Doel: Optimize the Critical Rendering Path   
[Exercises](https://github.com/cmda-minor-web/progressive-web-apps-2021/blob/master/course/week-3.md)  
[Critical Rendering Path - slides Declan Rek](https://github.com/cmda-minor-web/progressive-web-apps-1920/blob/master/course/cmd-2020-critical-rendering-path.pdf) -->


<!-- Add a nice image here at the end of the week, showing off your shiny frontend 📸 -->

<!-- What external data source is featured in your project and what are its properties 🌠 -->

<!-- How about a license here? 📜 (or is it a licence?) 🤷 -->
