import { createBrowserRouter, RouterProvider } from "react-router-dom"
import Base from "./components/layouts/Base"
import Projects from "./components/pages/Projects"


const Router= () => {
    const router = createBrowserRouter([
        {
            path: "/",
            element: <Base/>,
            children: [
                {
                    index: true,
                    element: <Projects />
                }
            ]
        },  
    ])

    return <RouterProvider router={router} />;
}

export default Router;