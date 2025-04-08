Here's a byte-level breakdown of a Git `receive-pack` protocol exchange over SSH, based on the provided references:

---

### **1. SSH Connection Setup**
- Client initiates SSH connection to server:  
  `ssh <user>@<host> "git-receive-pack '<repo>.git'"`   
  - The server authenticates the user and spawns the `git-receive-pack` process .

---

### **2. Server Advertises References (pkt-line format)**
Server sends references and capabilities using **pkt-line** encoding :
- Each pkt-line starts with a **4-byte hex length** (including the length field itself), followed by the payload.  
  Example:  
  ```
  0032776d796e636a756e6b000000000000 refs/heads/main\0 report-status delete-refs side-band-64k 
  ```
  - `0032` = 50 bytes (hex) → Total line length.  
  - `776d796e636a756e6b000000000000` = Object ID (SHA-1).  
  - `refs/heads/main` = Reference name.  
  - `\0` separates the ref from capabilities (e.g., `report-status`) .

---

### **3. Client Sends Commands (pkt-line format)**
Client sends push commands and packfile :
1. **Commands** (e.g., update/delete refs):  
   ```
   0028update refs/heads/main 776d796e636a756e6b000000000000 
   ```
   - `0028` = 40 bytes (hex).  
   - `update refs/heads/main` = Command to update the ref.  

2. **Packfile** (binary data):  
   - After sending commands, the client streams a **packfile** containing new objects .  
   - The packfile is prefixed with a `PACK` header and suffixed with a SHA-1 checksum .

---

### **4. Server Response (pkt-line format)**
Server sends status via **side-band-64k** channels :
- Example success response:  
  ```
  000eunpack ok 
  0018ok refs/heads/main 
  ```
  - `000e` = 14 bytes → `unpack ok` (packfile processed).  
  - `0018` = 24 bytes → `ok refs/heads/main` (ref updated).

---

### **Key Notes**
- **Authentication**: Handled entirely by SSH before `git-receive-pack` starts .  
- **No Negotiation**: Unlike fetch, the client sends a packfile directly since it already knows the server's refs .  
- **Push Options**: Sent as part of the initial command (e.g., `git push --push-option=...`) .

---

### **Byte Map Summary**
| Step | Direction | Data Format | Example Bytes |
|------|-----------|-------------|---------------|
| 1. SSH Handshake | Client → Server | SSH protocol | `ssh user@host...` |
| 2. Ref Advertisement | Server → Client | pkt-line | `0032... refs/heads/main\0...` |
| 3. Commands | Client → Server | pkt-line | `0028update refs/heads/main...` |
| 4. Packfile | Client → Server | PACK binary | `PACK...<objects>...<checksum>` |
| 5. Status | Server → Client | pkt-line | `000eunpack ok` |

---

This exchange ensures atomic updates to the repository  while relying on SSH for security .