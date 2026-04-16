import React from "react";

export interface CartInterface {
    // TODO
}

export default function CartComponent({ children }: React.PropsWithChildren<CartInterface>) {
    return (
        <Fragment>
            { children }
        </Fragment>
    )
}