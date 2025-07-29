import { createBrowserRouter, RouterProvider } from "react-router-dom"
import Base from "./components/layouts/Base"
import Projects from "./components/pages/Projects"
import ProjectDetails from "./components/pages/ProjectDetails"


const Router= () => {
    const router = createBrowserRouter([
        {
            path: "/",
            element: <Base/>,
            children: [
                {
                    index: true,
                    element: <Projects />
                },
                {
                    path: "/projects/:projectId",
                    element: <ProjectDetails />
                }
            ]
        },  
    ])

    return <RouterProvider router={router} />;
}

export default Router;