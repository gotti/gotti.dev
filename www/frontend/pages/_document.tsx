import Document, { Html, Head, Main, NextScript } from "next/document";
import { Menu } from "../components/Menu"

class MyDocument extends Document {
    render(): JSX.Element {
        return (
            <Html lang="ja">
                <Head />
                <body className="Body">
                    <div className="Menu"><Menu /></div>
                    <div className="Home"><Main /></div>
                    <NextScript />
                </body>
            </Html>
        );
    }
}

export default MyDocument;
