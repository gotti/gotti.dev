import { FC } from "react";
import { AppProps } from "next/app";
import { Menu } from "../components/Menu";

import "./app.scss";

const App: FC<AppProps> = ({ Component, pageProps }) => {
    return (
    <>
      <main>
        <link href="https://fonts.googleapis.com/css2?family=Source+Code+Pro:wght@300&family=Yomogi&display=swap" rel="stylesheet"/>
        <div className="Menu"><Menu/></div>
        <div className="Home"><Component {...pageProps} /></div>
      </main>
    </>
    );
};

export default App;
