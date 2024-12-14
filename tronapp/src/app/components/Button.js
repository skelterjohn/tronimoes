import { Button as AntButton } from 'antd';
import { forwardRef } from 'react';

const Button = forwardRef((props, ref) => {
  return <AntButton {...props} ref={ref} />;
});

Button.displayName = 'Button';

export default Button;
