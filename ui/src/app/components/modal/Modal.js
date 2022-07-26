

export const ModalHeader = ({children, onDismiss}) => {
  return (
    <div className="modal-header">
      <button type="button"
              className="close" 
              aria-label="Close"
              onClick={onDismiss}>
               <span aria-hidden="true">&times;</span></button>
      {children}
    </div>
  );
}


export const ModalBody = ({children}) => {
  return (
    <div className="modal-body">
      {children}
    </div>
  );
}


export const ModalFooter = ({children}) => {
  return(
    <div className="modal-footer">
      {children}
    </div>
  );
}


export const Modal = ({
  children,
  onDismiss,
  className="",
}) => {
  return (
    <div className={className}>
      <div className="modal modal-open modal-show fade in" role="dialog">
        <div className="modal-dialog" role="document">
          <div className="modal-content">
            {children}
          </div>
        </div>
      </div>
      <div className="modal-backdrop fade in"
           onClick={onDismiss}></div>
    </div>
  );
}

export default Modal;
