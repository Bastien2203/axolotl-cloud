import { useState } from "react";
import Button from "../atoms/Button";
import Modal from "../atoms/Modal";


const ImportComposeFileModal = ({ onClose, onImport }: { onClose: () => void; onImport: (file: string) => void; }) => {
    const [composeFile, setComposeFile] = useState("");
    return <Modal onClose={onClose}>
       <div className="min-w-[60vw] ">
            <h2>Import Compose File</h2>
            <br/>
            <textarea
                className="w-full h-64 p-2 border border-gray-300 rounded"
                placeholder="Paste your Docker Compose file here..."
                onChange={(e) => setComposeFile(e.target.value)}
            />
            <hr className="my-4 border-gray-300" />
            <Button onClick={() => onImport(composeFile)} className="w-full">Import</Button>
    </div>
    </Modal>
}

export default ImportComposeFileModal;