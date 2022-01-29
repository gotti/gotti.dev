import {FC} from "react";
import {AppProps} from "next/app";
import {Menu} from "../components/Menu";
import initTwitterScriptInner from 'zenn-embed-elements/lib/init-twitter-script-inner';
import "./app.scss";

const App: FC<AppProps> = ({Component, pageProps}) => {
  return (
    <>
      <script
        dangerouslySetInnerHTML={{
          __html: initTwitterScriptInner
        }}
      />
      <main>
        <div className="Menu">
          <Menu />
        </div>
        <div className="Home">
          <Component {...pageProps} />
        </div>
      </main>
    </>
  );
};

export default App;
