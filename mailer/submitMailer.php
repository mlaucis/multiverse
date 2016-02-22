<?php


require 'vendor/autoload.php';

use Aws\Ses\SesClient;

if (strtolower($_SERVER['REQUEST_METHOD']) == "options") {
    header("Access-Control-Allow-Origin: *");
    return;
}

$client = SesClient::factory(array(
    'key' => 'AKIAIHMKC34DIO4VSTSA',
    'secret' => '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',
    'region'  => 'eu-west-1',
));

$name = $_POST['name'];
$email = $_POST['email'];
$message = $_POST['message'];

if ($name == "" || $email == "" || $message == "") {
        http_response_code(400);
        return;
}

$fqea = $name.'<'.$email.'>';
$message = $fqea . ' wrote:<br/>'.$message;

$result = $client->sendEmail(array(
    'Source' => 'contact@tapglue.com',
    'Destination' => array(
        'ToAddresses' => array('Tapglue Contact <contact@tapglue.com>'),
        'BccAddresses' => array('Tapglue Contact <poic4idi@incoming.intercom.io>'),
    ),
    'Message' => array(
        'Subject' => array(
            'Data' => 'Website contact - ' . $name,
            'Charset' => 'UTF-8',
        ),
        'Body' => array(
            'Text' => array(
                'Data' => $message,
                'Charset' => 'UTF-8',
            ),
            'Html' => array(
                'Data' => $message,
                'Charset' => 'UTF-8',
            ),
        ),
    ),
    'ReplyToAddresses' => array($fqea),
    'ReturnPath' => 'bounce@tapglue.com',
));

header("Access-Control-Allow-Origin: *");
echo '{"status": "ok"}';
