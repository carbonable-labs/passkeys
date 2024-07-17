type ButtonProps = {
  onClick: () => void;
}

export function Button({ onClick, children }: React.PropsWithChildren<ButtonProps>) {
  return (
    <button onClick={onClick} className="mx-6 p-3 border rounded-md">{children}</button>
  )
}
