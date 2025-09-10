import json
import boto3
import urllib.parse
import base64
import os

lambda_client = boto3.client('lambda')

def lambda_handler(event, context):
    # Base64 인코딩된 body 디코딩 처리
    body = event.get('body', '')
    if event.get('isBase64Encoded', False):
        body = base64.b64decode(body).decode('utf-8')

    slack_data = urllib.parse.parse_qs(body)
    slack_dict = {k: v[0] if v else '' for k, v in slack_data.items()}

    command = slack_dict.get('text', '').strip()
    response_url = slack_dict.get('response_url')

    if command == 'list':
        blocks = [
            {"type": "section", "text": {"type": "plain_text", "text": "지원하는 서비스 및 명령어를 선택하세요 ✅"}},
            {"type": "actions", "elements": [
                {"type": "button", "text": {"type": "plain_text", "text": "s3 ls"}, "value": "s3 ls"},
                {"type": "button", "text": {"type": "plain_text", "text": "s3 mb"}, "value": "s3 mb"},
                {"type": "button", "text": {"type": "plain_text", "text": "s3 rb"}, "value": "s3 rb"},
                {"type": "button", "text": {"type": "plain_text", "text": "s3 cp"}, "value": "s3 cp"},
                {"type": "button", "text": {"type": "plain_text", "text": "s3 rm"}, "value": "s3 rm"},
                {"type": "button", "text": {"type": "plain_text", "text": "ec2 ... (준비중)"}, "value": "ec2"},
            ]}
        ]
        response = {
            "response_type": "in_channel",
            "blocks": blocks
        }
        return {
            "statusCode": 200,
            "headers": {"Content-Type": "application/json"},
            "body": json.dumps(response)
        }

    invoke_payload = {
        "command": command,
        "response_url": response_url,
        "user": slack_dict.get('user_name', 'unknown')
    }
    try:
        lambda_client.invoke(
            FunctionName=os.environ['ASYNC_LAMBDA_NAME'],
            InvocationType='Event',
            Payload=json.dumps(invoke_payload)
        )
    except Exception as e:
        err_text = f"비동기 처리 Lambda 호출 실패: {str(e)}"
        return {
            "statusCode": 200,
            "headers": {"Content-Type": "application/json"},
            "body": json.dumps({"text": err_text})
        }

    return {
        "statusCode": 200,
        "headers": {"Content-Type": "application/json"},
        "body": json.dumps({"text": "요청을 처리 중입니다... 잠시만 기다려 주세요."})
    }
