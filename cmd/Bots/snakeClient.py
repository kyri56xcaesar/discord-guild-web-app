import socket

HOST = "127.0.0.1"
PORT = 6970

def recv_until_newline(sock):
    """Receive data from the socket until a newline character is encountered."""
    data = b''
    while True:
        chunk = sock.recv(2048)
        if not chunk:
            return None
        data += chunk
        if b'\n' in chunk:
            break
    return data

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
    sock.connect((HOST, PORT))
    print(f"Connected to server at {HOST}:{PORT}")

    while True:
        inp = input("")
        if not inp:
            continue
        
        message = (inp + '\n').encode('utf-8')
        
        if len(message) > 2048:
            print("Error: Message too long.")
            continue
        sock.sendall(message)

        data = recv_until_newline(sock)
        if data is None:
            print("Server closed the connection.")
            break
        response = data.decode('utf-8').strip()
        print(response)
