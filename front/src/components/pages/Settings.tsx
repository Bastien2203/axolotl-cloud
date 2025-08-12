import { useEffect, useState } from "react";
import type { Setting } from "../../api/types";
import { getSettings } from "../../api/settings";
import { useToast } from "../../contexts/ToastContext";


const Settings = () => {
    const [settings, setSettings] = useState<Setting[]>();
    const toast = useToast();

    useEffect(() => {
        getSettings()
            .then(setSettings)
            .catch((error) => {
                console.error("Error fetching settings:", error);
                toast.error("Failed to fetch settings");
            });
    }, []);


    return <>
        <h1 className="text-2xl font-bold mb-4">Setting</h1>
        {
            settings ? (
            <div>
                {settings.map((setting) => (
                    <div key={setting.id} className="p-4 border-b border-gray-200 last:border-b-0">
                        <div className="flex justify-between items-center">
                            <span className="font-medium">{setting.key}</span>
                            <span className="text-gray-600">{setting.value}</span>
                        </div>
                    </div>
                ))}
            </div>
        ) : (<div className="text-gray-500">Loading settings...</div>)
        }
    </>
}

export default Settings;