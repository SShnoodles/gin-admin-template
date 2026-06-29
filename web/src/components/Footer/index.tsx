import { Divider } from 'antd';
import { createStyles } from 'antd-style';
import React from 'react';

const useStyles = createStyles(({ token, css }) => ({
  footer: css`
    padding: 16px 24px;
    text-align: center;
    color: ${token.colorTextDescription};
    font-size: ${token.fontSizeSM}px;
    line-height: ${token.lineHeight};
    background: transparent;
  `,
  copyright: css`
    margin-bottom: 6px;
  `,
  meta: css`
    display: flex;
    align-items: center;
    justify-content: center;
    flex-wrap: wrap;
    gap: 6px 12px;
    font-family: ${token.fontFamilyCode};
    font-size: ${token.fontSizeSM - 1}px;
  `,
  label: css`
    color: ${token.colorTextQuaternary};
  `,
  divider: css`
    display: inline-block;
    vertical-align: middle;
  `,
}));

const Footer: React.FC = () => {
  const { styles } = useStyles();
  const year = new Date().getFullYear();

  return (
    <div className={styles.footer}>
      <div className={styles.copyright}>Gin Admin &copy; {year}</div>
      <div className={styles.meta}>
        <span>
          <span className={styles.label}>ver </span>
          {__APP_VERSION__}
        </span>
        <Divider orientation="vertical" className={styles.divider} />
        <span>
          <span className={styles.label}>Umi </span>
          {__UMI_VERSION__}
        </span>
        <Divider orientation="vertical" className={styles.divider} />
        <span>
          <span className={styles.label}>Utoo </span>
          {__UTOO_VERSION__}
        </span>
      </div>
    </div>
  );
};

export default Footer;
