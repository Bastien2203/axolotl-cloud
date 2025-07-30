import { Folder, Settings } from "lucide-react";
import { Link, Outlet } from "react-router-dom"


const Base = () => {
    return <div className="flex">

        <div className="w-64 h-screen bg-gray-50 p-4">
            <Link className="flex flex-col items-center" to="/">
                <img src="/icon.png" alt="Axolotl Cloud Logo" className="w-16 h-16 mb-4" />
                <h1 className="text-xl">Axolotl Cloud</h1>
            </Link>

            <hr className="my-4 border-gray-300" />

            <nav>
                <NavItem to="/" label="Projects" icon={<Folder />} />
                <NavItem to="/settings" label="Settings" icon={<Settings />} />
            </nav>
        </div>
        <Outlet />
    </div>
}


type NavItemProps = {
    to: string;
    label: string;
    icon?: React.ReactNode;
}

const NavItem = (props: NavItemProps) => {
    return (
        <Link to={props.to} className="flex items-center p-2 text-gray-900 rounded-lg hover:bg-gray-100">
            {props.icon && <span className="mr-3">{props.icon}</span>}
            <span className="flex-1 ml-3 whitespace-nowrap">{props.label}</span>
        </Link>
    );
}

export default Base