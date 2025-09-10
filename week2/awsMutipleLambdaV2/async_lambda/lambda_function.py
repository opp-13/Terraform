import json
import boto3
import urllib3

http = urllib3.PoolManager()
s3 = boto3.client('s3')

def respond_slack(message, response_url):
    headers = {'Content-Type': 'application/json'}
    http.request('POST', response_url, body=json.dumps(message).encode('utf-8'), headers=headers)

def lambda_handler(event, context):
    command = event.get('command', '')
    response_url = event.get('response_url', '')
    user = event.get('user', 'unknown')

    text = f"{user}님의 명령어 처리 결과:\n"

    try:
        if command.startswith('s3 ls'):
            buckets = s3.list_buckets()
            bucket_list = "\n".join([b['Name'] for b in buckets.get('Buckets', [])])
            text += bucket_list or '버킷이 없습니다.'
        elif command.startswith('s3 mb'):
            bucket_name = command.split()[-1].replace('s3://', '')
            s3.create_bucket(Bucket=bucket_name)
            text += f"버킷 생성 완료: {bucket_name}"
        elif command.startswith('s3 rb'):
            bucket_name = command.split()[-1].replace('s3://', '')
            s3_resource = boto3.resource('s3')
            bucket = s3_resource.Bucket(bucket_name)
            bucket.object_versions.delete()
            bucket.delete()
            text += f"버킷 삭제 완료: {bucket_name}"
        elif command.startswith('s3 cp'):
            tokens = command.split()
            src = tokens[2].replace('s3://', '')
            dest = tokens[3].replace('s3://', '')
            src_bucket, src_key = src.split('/', 1)
            dest_bucket, dest_key = dest.split('/', 1)
            copy_source = {'Bucket': src_bucket, 'Key': src_key}
            s3.copy_object(Bucket=dest_bucket, Key=dest_key, CopySource=copy_source)
            text += f"객체 복사 완료: {src_bucket}/{src_key} → {dest_bucket}/{dest_key}"
        elif command.startswith('s3 rm'):
            obj = command.split()[-1].replace('s3://', '')
            bucket_name, key_name = obj.split('/', 1)
            s3.delete_object(Bucket=bucket_name, Key=key_name)
            text += f"객체 삭제 완료: {bucket_name}/{key_name}"
        else:
            text += "지원하지 않는 명령어입니다."
    except Exception as e:
        text += f"명령 처리 중 오류 발생: {str(e)}"

    if response_url:
        respond_slack({"text": text}, response_url)
