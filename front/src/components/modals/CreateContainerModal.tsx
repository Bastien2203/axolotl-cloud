import { useState } from "react";
import { useToast } from "../../contexts/ToastContext";
import Button from "../atoms/Button";
import Input from "../atoms/Input";
import Modal from "../atoms/Modal";
import type { Container, NetworkMode } from "../../api/types";
import KeyValueEditor from "../atoms/KeyValueEditor";
import StringListEditor from "../atoms/StringListEditor";
import Select from "../atoms/Select";


const CreateContainerModal = ({ onClose, onCreate, defaultValue }: { onClose: () => void, onCreate: (container: Omit<Container, 'id'>) => void, defaultValue?: Omit<Container, 'id'> }) => {
    const [newContainer, setNewContainer] = useState<Omit<Container, 'id'>>(defaultValue || {
        name: "",
        docker_image: "",
        ports: {},
        env: {},
        volumes: {},
        networks: [],
        network_mode: "bridge", // Default network mode
    });

    const toast = useToast();

    const validateForm = (container: Omit<Container, 'id'>) => {
        if (!container.name || !container.docker_image) {
            return false;
        }
        return true;
    }

    return <Modal onClose={onClose}>
        <div className="min-w-[60vw]">
            <h2>Create New Container</h2>
            <br />
            <form onSubmit={(e) => {
                e.preventDefault();
                if (validateForm(newContainer)) {
                    onCreate(newContainer)
                } else {
                    toast.error("Please fill in all required fields.");
                }
            }} className="flex flex-col gap-4">
                <h3>Container Details</h3>
                <Input type="text" name="name" className="w-full " placeholder="Container name" required onChange={(e) => setNewContainer({ ...newContainer, name: e.target.value })} value={newContainer.name} />
                <Input type="text" name="docker_image" className="w-full" placeholder="Docker Image" required onChange={(e) => setNewContainer({ ...newContainer, docker_image: e.target.value })} value={newContainer.docker_image} />

                <div className="flex items-center gap-2">
                <label htmlFor="network_mode" className="block text-sm font-medium text-gray-700 flex-shrink-0">
                    Network Mode
                </label>
                <Select
                    name="network_mode"
                    className="w-full"
                    value={newContainer.network_mode}
                    onChange={(e) => setNewContainer({ ...newContainer, network_mode: e.target.value as NetworkMode })}
                    options={[
                        { label: "Bridge", value: "bridge" },
                        { label: "Host", value: "host" },
                        { label: "None", value: "none" },
                    ]}
                />
                </div>

                <KeyValueEditor
                    label="Ports"
                    addLabel="Add Port"
                    data={newContainer.ports || {}}
                    onChange={(ports) => setNewContainer({ ...newContainer, ports })}
                    placeholderKey="Host Port"
                    placeholderValue="Container Port"
                    variant="secondary"
                />

                <KeyValueEditor
                    label="Environment Variables"
                    addLabel="Add Variable"
                    data={newContainer.env || {}}
                    onChange={(env) => setNewContainer({ ...newContainer, env })}
                    placeholderKey="Key"
                    placeholderValue="Value"
                    variant="secondary"
                />

                <KeyValueEditor
                    label="Volumes"
                    addLabel="Add Volume"
                    data={newContainer.volumes || {}}
                    onChange={(volumes) => setNewContainer({ ...newContainer, volumes })}
                    placeholderKey="Host Path"
                    placeholderValue="Container Path"
                    variant="secondary"
                />

                <StringListEditor
                    label="Networks"
                    data={newContainer.networks || []}
                    onChange={(networks) => setNewContainer({ ...newContainer, networks })}
                    addLabel="Add Network"
                    placeholder="network_name"
                    variant="secondary"
                />



                <hr className="my-4 border-gray-300" />

                <Button type="submit">{defaultValue ? "Save Changes" : "Create Container"}</Button>
            </form>

        </div>
    </Modal>
}

export default CreateContainerModal;