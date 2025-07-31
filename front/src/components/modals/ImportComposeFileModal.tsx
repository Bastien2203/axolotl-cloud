import { useState } from "react";
import Button from "../atoms/Button";
import Modal from "../atoms/Modal";
import Editor from "@monaco-editor/react";

const ImportComposeFileModal = ({
  onClose,
  onImport,
}: {
  onClose: () => void;
  onImport: (file: string) => void;
}) => {
  const [composeFile, setComposeFile] = useState("");

  return (
    <Modal onClose={onClose}>
      <div className="min-w-[60vw]">
        <h2 className="text-xl font-semibold mb-4">Import Compose File</h2>
        <div className="h-64 border border-gray-300 rounded overflow-hidden">
          <Editor
            height="100%"
            defaultLanguage="yaml"
            defaultValue="# Paste your Docker Compose file here..."
            value={composeFile}
            onChange={(value) => setComposeFile(value || "")}
            theme="vs-light"
            options={{
              minimap: { enabled: false },
              fontSize: 14,
              automaticLayout: true,
            }}
          />
        </div>
        <hr className="my-4 border-gray-300" />
        <Button onClick={() => onImport(composeFile)} className="w-full">
          Import
        </Button>
      </div>
    </Modal>
  );
};

export default ImportComposeFileModal;
