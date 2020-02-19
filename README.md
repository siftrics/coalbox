# Coalbox

The [Sight API](https://siftrics.com/) is a text recognition service. It provides word-level bounding boxes in response to uploaded PDF documents and images.

This repository provides a Go library of bounding box coalescence algorithms: it is a toolkit of functions which take word-level bounding boxes as input and return coalesced &mdash; e.g., sentence-level or paragraph-level &mdash; bounding boxes.
