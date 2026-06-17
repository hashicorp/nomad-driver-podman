<!doctype html>
<html>
  <head>
    <title>{{ env "SITE_NAME" }}</title>
  </head>
  <body>
    <h1>Hello from {{ env "SITE_NAME" }}!</h1>
    <p>Served by allocation <code>{{ env "NOMAD_ALLOC_ID" }}</code></p>
    <p>Running on node <code>{{ env "node.unique.name" }}</code></p>
  </body>
</html>
