import { Trash, Plus } from "lucide-react";
import Button from "./Button";
import Input from "./Input";


type KeyValueEditorProps = {
  label: string;
  data: Record<string, string>;
  onChange: (updated: Record<string, string>) => void;
  addLabel: string;
  variant?: "primary" | "secondary" | "danger";
  placeholderKey?: string;
  placeholderValue?: string;
};

const KeyValueEditor = ({ label, data, onChange, addLabel, variant, placeholderKey, placeholderValue }: KeyValueEditorProps) => {
  return (
    <div>
      <h3>{label}</h3>
      <div className="flex flex-col gap-2">
        {Object.entries(data).map(([k, v], i) => (
          <div className="flex gap-2 items-center" key={`${label}-${i}`}>
            <Input
              name={`${label}_${i}_key`}
              value={k}
              placeholder={placeholderKey}
              onChange={(e) => {
                const entries = Object.entries(data);
                entries[i][0] = e.target.value;
                const updated = Object.fromEntries(entries);
                onChange(updated);
              }}
              className="flex-1"
            />
            <Input
              name={`${label}_${i}_value`}
              value={v}
              placeholder={placeholderValue}
              onChange={(e) => {
                const entries = Object.entries(data);
                entries[i][1] = e.target.value;
                const updated = Object.fromEntries(entries);
                onChange(updated);
              }}
              className="flex-1"
            />
            <Trash
              className="text-red-500 cursor-pointer flex-shrink-0"
              onClick={() => {
                const entries = Object.entries(data);
                entries.splice(i, 1);
                const updated = Object.fromEntries(entries);
                onChange(updated);
              }}
            />
          </div>
        ))}
      </div>
      <Button
        onClick={() => onChange({ ...data, "": "" })}
        variant={variant}
        className="my-2 gap-2 w-full"
      >
        <Plus />
        {addLabel}
      </Button>
    </div>
  );
};

export default KeyValueEditor;