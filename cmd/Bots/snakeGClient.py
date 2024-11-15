import socket
import tkinter as tk
from tkinter import simpledialog, scrolledtext

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

def run_gui(sock):
    def send_message(event=None):
        message = entry.get("1.0", 'end-1c') + '\n'
        if not message.strip():
            return
        if len(message.encode('utf-8')) > 2048:
            text_area.config(state=tk.NORMAL)
            text_area.insert(tk.END, "Error: Message too long.\n")
            text_area.config(state=tk.DISABLED)
            return
        sock.sendall(message.encode('utf-8'))
        entry.delete("1.0", tk.END)

        data = recv_until_newline(sock)
        if data is None:
            text_area.config(state=tk.NORMAL)
            text_area.insert(tk.END, "Server closed the connection.\n")
            text_area.config(state=tk.DISABLED)
            entry.config(state=tk.DISABLED)
            send_button.config(state=tk.DISABLED)
            return
        response = data.decode('utf-8').strip()
        text_area.config(state=tk.NORMAL)
        text_area.insert(tk.END, f"Server: {response}\n")
        text_area.config(state=tk.DISABLED)

    root = tk.Tk()
    root.title("Snake Bot Client")
    root.configure(bg="#2e2e2e")

    text_area = scrolledtext.ScrolledText(root, wrap=tk.WORD, state='normal', height=20, width=50, bg="#1e1e1e", fg="#dcdcdc", insertbackground="white")
    text_area.pack(padx=10, pady=5)
    text_area.insert(tk.END, f"Connected to server at {HOST}:{PORT}\n")
    text_area.config(state=tk.DISABLED)

    entry = tk.Text(root, height=3, width=50, bg="#1e1e1e", fg="#dcdcdc", insertbackground="white")
    entry.pack(padx=10, pady=5)
    entry.bind("<Return>", send_message)

    send_button = tk.Button(root, text="Send", command=send_message, bg="#3e3e3e", fg="#dcdcdc")
    send_button.pack(padx=10, pady=5)

    root.mainloop()

# Command-line interface
with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
    sock.connect((HOST, PORT))
    print(f"Connected to server at {HOST}:{PORT}")

    import sys
    if len(sys.argv) > 1 and sys.argv[1] == '-g':
        run_gui(sock)
    else:
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
