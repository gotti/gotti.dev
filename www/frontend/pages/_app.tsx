import { FC } from "react";
import { AppProps } from "next/app";

import "./app.scss"

const App: FC<AppProps> = ({ Component, pageProps }) => {
    return <Component {...pageProps} />;
};

export default App;
