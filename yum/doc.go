/*
A native implementation of Yum package repositories in Go.

This package contains functions to interrogate package repository metadata and
databases that were built using the createrepo utility.

It deliberately avoids reimplementing functionality of the yum tool itself (such
as the downloading and caching of repository metadata and databases). Some of
these functions are timely and should be implemented according to target user
interface (E.g. with progress bars and configurable thread counts). Other
functions of the yum tool are geared toward manipulation of the target system
(such as the installation and removal of packages) and are also outside the
scope of this package.

*/
package yum
