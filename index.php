<?php

ini_set('default_socket_timeout', $_ENV['PHP_SOCKET_TIMEOUT'] ?? 10);
ini_set('display_errors', 0);

header('Content-Type:text/plain');

$url = sprintf('%s:%s/xml/modified.xml?password=%s', $_ENV['CORE_HOST'], $_ENV['CORE_PORT'], $_ENV['CORE_PASSWORD']);
$response = file_get_contents($url);

if (false === $response) {
    header('HTTP/1.0 500 Internal Server Error');
    printf('no response from %s', $url);
    die;
}

$xml = new SimpleXMLElement($response);

foreach ($xml->information->attributes() as $key => $value) {
    printf('applejuice_%s %s%s', $key, $value, PHP_EOL);
}
