# This is the generated MFP model file.
# You probably need to edit it appropriately before use.

# IPP printer attributes:
ipp.attrs = {
  'operations-supported': [
    ipp.OP.PRINT_JOB,
    ipp.OP.VALIDATE_JOB,
    ipp.OP.CANCEL_JOB,
    ipp.OP.SEND_DOCUMENT,
    ipp.OP.CREATE_JOB,
    ipp.OP.GET_JOB_ATTRIBUTES,
    ipp.OP.GET_JOBS,
    ipp.OP.GET_PRINTER_ATTRIBUTES,
    ipp.OP.IDENTIFY_PRINTER
  ],
  'ipp-versions-supported': [
    ipp.KEYWORD('1.0'),
    ipp.KEYWORD('1.1'),
    ipp.KEYWORD('2.0')
  ],
  'charset-configured': ipp.CHARSET('utf-8'),
  'charset-supported': ipp.CHARSET('utf-8'),
  'natural-language-configured': ipp.LANGUAGE('en-us'),
  'generated-natural-language-supported': ipp.LANGUAGE('en-us'),
  'document-format-default': ipp.MIMETYPE('application/octet-stream'),
  'pdl-override-supported': ipp.KEYWORD('attempted'),
  'compression-supported': ipp.KEYWORD('none'),
  'multiple-document-jobs-supported': ipp.BOOLEAN(False),
  'multiple-operation-time-out': ipp.INTEGER(120),
  'media-col-supported': [
    ipp.KEYWORD('media-size'),
    ipp.KEYWORD('media-bottom-margin'),
    ipp.KEYWORD('media-left-margin'),
    ipp.KEYWORD('media-right-margin'),
    ipp.KEYWORD('media-top-margin'),
    ipp.KEYWORD('media-source'),
    ipp.KEYWORD('media-type')
  ],
  'finishings-default': ipp.ENUM(3),
  'pdf-versions-supported': ipp.KEYWORD('iso-32000-1_2008'),
  'landscape-orientation-requested-preferred': ipp.ENUM(4),
  'orientation-requested-default': ipp.ENUM(3),
  'orientation-requested-supported': [
    ipp.ENUM(3),
    ipp.ENUM(4),
    ipp.ENUM(5),
    ipp.ENUM(6),
    ipp.ENUM(7)
  ],
  'job-creation-attributes-supported': [
    ipp.KEYWORD('orientation-requested'),
    ipp.KEYWORD('copies'),
    ipp.KEYWORD('finishings'),
    ipp.KEYWORD('ipp-attribute-fidelity'),
    ipp.KEYWORD('job-name'),
    ipp.KEYWORD('media'),
    ipp.KEYWORD('media-col'),
    ipp.KEYWORD('operation-requested'),
    ipp.KEYWORD('output-bin'),
    ipp.KEYWORD('print-quality'),
    ipp.KEYWORD('printer-resolution'),
    ipp.KEYWORD('sides'),
    ipp.KEYWORD('print-color-mode'),
    ipp.KEYWORD('multiple-document-handling'),
    ipp.KEYWORD('page-ranges'),
    ipp.KEYWORD('page-content-optimize'),
    ipp.KEYWORD('page-scaling'),
    ipp.KEYWORD('feed-orientation'),
    ipp.KEYWORD('overrides'),
    ipp.KEYWORD('job-mandatory-attributes')
  ],
  'media-bottom-margin-supported': ipp.INTEGER(400),
  'media-left-margin-supported': ipp.INTEGER(400),
  'media-right-margin-supported': ipp.INTEGER(400),
  'media-top-margin-supported': ipp.INTEGER(400),
  'media-type-supported': [
    ipp.KEYWORD('auto'),
    ipp.KEYWORD('stationery'),
    ipp.KEYWORD('transparency'),
    ipp.KEYWORD('envelope'),
    ipp.KEYWORD('labels'),
    ipp.KEYWORD('cardstock'),
    ipp.KEYWORD('stationery-lightweight'),
    ipp.KEYWORD('stationery-preprinted'),
    ipp.KEYWORD('stationery-bond'),
    ipp.KEYWORD('stationery-colored'),
    ipp.KEYWORD('stationery-prepunched'),
    ipp.KEYWORD('stationery-letterhead'),
    ipp.KEYWORD('stationery-heavyweight'),
    ipp.KEYWORD('stationery-fine')
  ],
  'page-ranges-supported': ipp.BOOLEAN(True),
  'printer-resolution-default': ipp.RESOLUTION(600, 600, 'dpi'),
  'printer-resolution-supported': ipp.RESOLUTION(600, 600, 'dpi'),
  'print-quality-default': ipp.ENUM(4),
  'print-quality-supported': ipp.ENUM(4),
  'job-priority-default': ipp.INTEGER(50),
  'job-priority-supported': ipp.INTEGER(100),
  'identify-actions-default': ipp.KEYWORD('sound'),
  'identify-actions-supported': [ipp.KEYWORD('flash'), ipp.KEYWORD('sound')],
  'print-content-optimize-default': ipp.KEYWORD('auto'),
  'print-content-optimize-supported': ipp.KEYWORD('auto'),
  'print-scaling-supported': [
    ipp.KEYWORD('auto'),
    ipp.KEYWORD('auto-fit'),
    ipp.KEYWORD('fill'),
    ipp.KEYWORD('fit'),
    ipp.KEYWORD('none')
  ],
  'printer-kind': [ipp.KEYWORD('document'), ipp.KEYWORD('envelope')],
  'job-ids-supported': ipp.BOOLEAN(True),
  'which-jobs-supported': [ipp.KEYWORD('completed'), ipp.KEYWORD('not-completed')],
  'multiple-operation-time-out-action': ipp.KEYWORD('abort-job'),
  'document-format-version-supplied': ipp.KEYWORD('1.4'),
  'printer-organization': ipp.TEXT(''),
  'printer-organizational-unit': ipp.TEXT(''),
  'feed-orientation-supported': ipp.KEYWORD('short-edge-first'),
  'pwg-raster-document-sheet-back': ipp.KEYWORD('normal'),
  'pwg-raster-document-resolution-supported': ipp.RESOLUTION(600, 600, 'dpi'),
  'pwg-raster-document-type-supported': [
    ipp.KEYWORD('black_1'),
    ipp.KEYWORD('sgray_8'),
    ipp.KEYWORD('srgb_8'),
    ipp.KEYWORD('cmyk_8')
  ],
  'print_wfds': ipp.TEXT('T'),
  'mopria-certified': ipp.TEXT('1.3'),
  'overrides-supported': [
    ipp.KEYWORD('pages'),
    ipp.KEYWORD('document-numbers'),
    ipp.KEYWORD('media'),
    ipp.KEYWORD('media-col')
  ],
  'job-pages-per-set-supported': ipp.BOOLEAN(False),
  'jpeg-features-supported': [ipp.KEYWORD('none'), ipp.KEYWORD('cmyk')],
  'copies-supported': ipp.RANGE(1, 999),
  'printer-uri-supported': [
    ipp.URI('ipps://192.168.1.102:443/ipp/print'),
    ipp.URI('ipp://192.168.1.102:631/ipp/print')
  ],
  'media-supported': [
    ipp.KEYWORD('iso_a4_210x297mm'),
    ipp.KEYWORD('iso_a5_148x210mm'),
    ipp.KEYWORD('iso_a6_105x148mm'),
    ipp.KEYWORD('iso_b5_176x250mm'),
    ipp.KEYWORD('na_legal_8.5x14in'),
    ipp.KEYWORD('na_letter_8.5x11in'),
    ipp.KEYWORD('na_executive_7.25x10.5in'),
    ipp.KEYWORD('na_invoice_5.5x8.5in'),
    ipp.KEYWORD('iso_c5_162x229mm'),
    ipp.KEYWORD('iso_c6_114x162mm'),
    ipp.KEYWORD('iso_dl_110x220mm'),
    ipp.KEYWORD('na_monarch_3.875x7.5in'),
    ipp.KEYWORD('jis_b5_182x257mm'),
    ipp.KEYWORD('jis_b6_128x182mm'),
    ipp.KEYWORD('jpn_you4_105x235mm'),
    ipp.KEYWORD('jpn_hagaki_100x148mm'),
    ipp.KEYWORD('jpn_oufuku_148x200mm'),
    ipp.KEYWORD('roc_16k_7.75x10.75in'),
    ipp.KEYWORD('na_foolscap_8.5x13in'),
    ipp.KEYWORD('na_number-10_4.125x9.5in'),
    ipp.KEYWORD('na_number-9_3.875x8.875in'),
    ipp.KEYWORD('na_personal_3.625x6.5in'),
    ipp.KEYWORD('om_folio_210x330mm'),
    ipp.KEYWORD('custom_max_216x356mm'),
    ipp.KEYWORD('custom_min_70x148mm')
  ],
  'jpeg-x-dimension-supported': ipp.RANGE(1, 6000),
  'jpeg-y-dimension-supported': ipp.RANGE(1, 8410),
  'uri-security-supported': [ipp.KEYWORD('tls'), ipp.KEYWORD('none')],
  'printer-name': ipp.NAME('KM7B6A91'),
  'printer-location': ipp.TEXT('Living Room'),
  'printer-make-and-model': ipp.TEXT('ECOSYS M2040dn'),
  'color-supported': ipp.BOOLEAN(False),
  'print-color-mode-supported': [
    ipp.KEYWORD('monochrome'),
    ipp.KEYWORD('auto'),
    ipp.KEYWORD('auto-monochrome')
  ],
  'sides-supported': [
    ipp.KEYWORD('one-sided'),
    ipp.KEYWORD('two-sided-short-edge'),
    ipp.KEYWORD('two-sided-long-edge')
  ],
  'finishings-supported': ipp.ENUM(3),
  'output-bin-supported': [ipp.KEYWORD('top'), ipp.KEYWORD('face-down')],
  'media-source-supported': [
    ipp.KEYWORD('auto'),
    ipp.KEYWORD('by-pass-tray'),
    ipp.KEYWORD('tray-1')
  ],
  'jpeg-k-octets-supported': ipp.RANGE(0, 49152),
  'pdf-k-octets-supported': ipp.RANGE(0, 49152),
  'multiple-document-handling-supported': [
    ipp.KEYWORD('separate-documents-collated-copies'),
    ipp.KEYWORD('separate-documents-uncollated-copies')
  ],
  'document-format-supported': [
    ipp.MIMETYPE('application/octet-stream'),
    ipp.MIMETYPE('application/pdf'),
    ipp.MIMETYPE('image/tiff'),
    ipp.MIMETYPE('image/jpeg'),
    ipp.MIMETYPE('image/urf'),
    ipp.MIMETYPE('application/postscript'),
    ipp.MIMETYPE('application/vnd.hp-PCL'),
    ipp.MIMETYPE('application/vnd.hp-PCLXL'),
    ipp.MIMETYPE('application/vnd.xpsdocument'),
    ipp.MIMETYPE('image/pwg-raster')
  ],
  'uri-authentication-supported': [ipp.KEYWORD('none'), ipp.KEYWORD('none')],
  'media-default': ipp.KEYWORD('iso_a4_210x297mm'),
  'copies-default': ipp.INTEGER(1),
  'output-bin-default': ipp.KEYWORD('top'),
  'print-color-mode-default': ipp.KEYWORD('monochrome'),
  'sides-default': ipp.KEYWORD('one-sided'),
  'print-scaling-default': ipp.KEYWORD('auto'),
  'feed-orientation-default': ipp.KEYWORD('short-edge-first'),
  'pages-per-minute': ipp.INTEGER(40),
  'printer-uuid': ipp.URI('urn:uuid:4509a320-00a0-008f-00b6-002507510eca'),
  'printer-device-id': ipp.TEXT('ID:ECOSYS M2040dn;MFG:Kyocera;CMD:PCLXL,PostScript Emulation,PCL5E,PJL;MDL:ECOSYS M2040dn;CLS:PRINTER;DES:Kyocera ECOSYS M2040dn;SER:VCF9192281;URF:CP255,DM4,IFU0,IS19-20,OB1-10,PQ4,RS600,V1.4,W8;'),
  'printer-geo-location': ipp.UNKNOWN(),
  'marker-colors': [ipp.NAME('#000000'), ipp.NAME('none')],
  'marker-high-levels': [ipp.INTEGER(100), ipp.INTEGER(95)],
  'marker-levels': [ipp.INTEGER(35), ipp.INTEGER(0)],
  'marker-low-levels': [ipp.INTEGER(5), ipp.INTEGER(0)],
  'marker-names': [ipp.NAME('Black TK-1170'), ipp.NAME('Waste Toner Box')],
  'marker-types': [ipp.KEYWORD('toner'), ipp.KEYWORD('waste-toner')],
  'printer-dns-sd-name': ipp.NAME('Kyocera ECOSYS M2040dn'),
  'printer-info': ipp.TEXT('Kyocera ECOSYS M2040dn'),
  'printer-alert': ipp.STRING('other'),
  'printer-alert-description': ipp.TEXT('Sleeping... '),
  'printer-input-tray': [
    ipp.STRING('type=other;mediafeed=116929;mediaxfeed=82677;mediafeed=116929;mediaxfeed=82677;maxcapacity=-2;level=-2;status=19;name=Auto;'),
    ipp.STRING('type=sheetFeedAutoNonRemovableTray;mediafeed=116929;mediaxfeed=82677;maxcapacity=100;level=0;status=19;name=MP Tray;'),
    ipp.STRING('type=sheetFeedAutoNonRemovableTray;mediafeed=116929;mediaxfeed=82677;maxcapacity=250;level=-3;status=19;name=Cassette 1;')
  ],
  'printer-output-tray': [
    ipp.STRING('type=unRemovableBin;maxcapacity=150;remaining=-2;status=4;name=Top Tray;stackingorder=firstToLast;pagedelivery=faceDown;'),
    ipp.STRING('type=unRemovableBin;maxcapacity=150;remaining=-2;status=4;name=Top Tray;stackingorder=firstToLast;pagedelivery=faceDown;')
  ],
  'printer-firmware-name': ipp.NAME('System'),
  'printer-firmware-string-version': ipp.TEXT('2S0_2000.001.828'),
  'printer-firmware-version': ipp.STRING('2S0_2000.001.828'),
  'multiple-document-handling-default': ipp.KEYWORD('separate-documents-collated-copies'),
  'printer-state-message': ipp.TEXT('Sleeping...  '),
  'urf-supported': [
    ipp.KEYWORD('CP255'),
    ipp.KEYWORD('DM4'),
    ipp.KEYWORD('IFU0'),
    ipp.KEYWORD('IS19-20'),
    ipp.KEYWORD('OB1-10'),
    ipp.KEYWORD('PQ4'),
    ipp.KEYWORD('RS600'),
    ipp.KEYWORD('V1.4'),
    ipp.KEYWORD('W8')
  ],
  'printer-up-time': ipp.INTEGER(3531253),
  'queued-job-count': ipp.INTEGER(0),
  'printer-is-accepting-jobs': ipp.BOOLEAN(True),
  'printer-state': ipp.ENUM(3),
  'printer-state-reasons': ipp.KEYWORD('none'),
  'ipp-features-supported': [
    ipp.KEYWORD('airprint-1.3'),
    ipp.KEYWORD('airprint-1.4'),
    ipp.KEYWORD('airprint-1.5'),
    ipp.KEYWORD('airprint-1.6')
  ],
  'printer-more-info': ipp.URI('https://192.168.1.102/airprint'),
  'printer-supply-info-uri': ipp.URI('https://192.168.1.102'),
  'printer-fax-log-uri': ipp.URI('https://192.168.1.102/printer-fax-log/faxout.log'),
  'printer-icons': [
    ipp.URI('https://192.168.1.102/printer-icon/machine_128.png'),
    ipp.URI('https://192.168.1.102/printer-icon/machine_512.png')
  ],
  'media-col-default': {
    'media-type': ipp.KEYWORD('stationery'),
    'media-size': {'x-dimension': ipp.INTEGER(21000), 'y-dimension': ipp.INTEGER(29700)},
    'media-top-margin': ipp.INTEGER(400),
    'media-left-margin': ipp.INTEGER(400),
    'media-right-margin': ipp.INTEGER(400),
    'media-bottom-margin': ipp.INTEGER(400),
    'media-source': ipp.KEYWORD('by-pass-tray')
  },
  'media-size-supported': [
    {'x-dimension': ipp.INTEGER(21000), 'y-dimension': ipp.INTEGER(29700)},
    {'x-dimension': ipp.INTEGER(14800), 'y-dimension': ipp.INTEGER(21000)},
    {'x-dimension': ipp.INTEGER(10500), 'y-dimension': ipp.INTEGER(14800)},
    {'x-dimension': ipp.INTEGER(17600), 'y-dimension': ipp.INTEGER(25000)},
    {'x-dimension': ipp.INTEGER(21590), 'y-dimension': ipp.INTEGER(35560)},
    {'x-dimension': ipp.INTEGER(21590), 'y-dimension': ipp.INTEGER(27940)},
    {'x-dimension': ipp.INTEGER(18415), 'y-dimension': ipp.INTEGER(26670)},
    {'x-dimension': ipp.INTEGER(13970), 'y-dimension': ipp.INTEGER(21590)},
    {'x-dimension': ipp.INTEGER(16200), 'y-dimension': ipp.INTEGER(22900)},
    {'x-dimension': ipp.INTEGER(11400), 'y-dimension': ipp.INTEGER(16200)},
    {'x-dimension': ipp.INTEGER(11000), 'y-dimension': ipp.INTEGER(22000)},
    {'x-dimension': ipp.INTEGER(9842), 'y-dimension': ipp.INTEGER(19050)},
    {'x-dimension': ipp.INTEGER(18200), 'y-dimension': ipp.INTEGER(25700)},
    {'x-dimension': ipp.INTEGER(12800), 'y-dimension': ipp.INTEGER(18200)},
    {'x-dimension': ipp.INTEGER(10500), 'y-dimension': ipp.INTEGER(23500)},
    {'x-dimension': ipp.INTEGER(10000), 'y-dimension': ipp.INTEGER(14800)},
    {'x-dimension': ipp.INTEGER(14800), 'y-dimension': ipp.INTEGER(20000)},
    {'x-dimension': ipp.INTEGER(19685), 'y-dimension': ipp.INTEGER(27305)},
    {'x-dimension': ipp.INTEGER(21590), 'y-dimension': ipp.INTEGER(33020)},
    {'x-dimension': ipp.INTEGER(10477), 'y-dimension': ipp.INTEGER(24130)},
    {'x-dimension': ipp.INTEGER(9842), 'y-dimension': ipp.INTEGER(22542)},
    {'x-dimension': ipp.INTEGER(9207), 'y-dimension': ipp.INTEGER(16510)},
    {'x-dimension': ipp.INTEGER(21000), 'y-dimension': ipp.INTEGER(33000)},
    {
      'x-dimension': ipp.RANGE(7000, 21600),
      'y-dimension': ipp.RANGE(14800, 35600)
    }
  ],
  'media-col-ready': {
    'media-type': ipp.KEYWORD('stationery'),
    'media-size': {'x-dimension': ipp.INTEGER(21000), 'y-dimension': ipp.INTEGER(29700)},
    'media-top-margin': ipp.INTEGER(400),
    'media-left-margin': ipp.INTEGER(400),
    'media-right-margin': ipp.INTEGER(400),
    'media-bottom-margin': ipp.INTEGER(400),
    'media-source': ipp.KEYWORD('tray-1'),
    'media-source-properties': {'media-source-feed-direction': ipp.KEYWORD('short-edge-first'), 'media-source-feed-orientation': ipp.ENUM(3)}
  },
  'media-ready': ipp.KEYWORD('iso_a4_210x297mm'),
  'printer-get-attributes-supported': [
    ipp.KEYWORD('operations-supported'),
    ipp.KEYWORD('ipp-versions-supported'),
    ipp.KEYWORD('charset-configured'),
    ipp.KEYWORD('charset-supported'),
    ipp.KEYWORD('natural-language-configured'),
    ipp.KEYWORD('generated-natural-language-supported'),
    ipp.KEYWORD('document-format-default'),
    ipp.KEYWORD('document-format-supported'),
    ipp.KEYWORD('pdl-override-supported'),
    ipp.KEYWORD('compression-supported'),
    ipp.KEYWORD('multiple-document-jobs-supported'),
    ipp.KEYWORD('multiple-operation-time-out'),
    ipp.KEYWORD('multiple-document-handling-default'),
    ipp.KEYWORD('multiple-document-handling-supported'),
    ipp.KEYWORD('copies-supported'),
    ipp.KEYWORD('media-col-supported'),
    ipp.KEYWORD('pdf-versions-supported'),
    ipp.KEYWORD('orientation-requested-default'),
    ipp.KEYWORD('orientation-requested-supported'),
    ipp.KEYWORD('job-creation-attributes-supported'),
    ipp.KEYWORD('media-bottom-margin-supported'),
    ipp.KEYWORD('media-left-margin-supported'),
    ipp.KEYWORD('media-right-margin-supported'),
    ipp.KEYWORD('media-top-margin-supported'),
    ipp.KEYWORD('page-ranges-supported'),
    ipp.KEYWORD('printer-resolution-default'),
    ipp.KEYWORD('printer-resolution-supported'),
    ipp.KEYWORD('print-quality-default'),
    ipp.KEYWORD('print-quality-supported'),
    ipp.KEYWORD('job-priority-default'),
    ipp.KEYWORD('job-priority-supported'),
    ipp.KEYWORD('identify-actions-default'),
    ipp.KEYWORD('identify-actions-supported'),
    ipp.KEYWORD('print-content-optimize-default'),
    ipp.KEYWORD('print-content-optimize-supported'),
    ipp.KEYWORD('print-scaling-supported'),
    ipp.KEYWORD('ipp-features-supported'),
    ipp.KEYWORD('job-ids-supported'),
    ipp.KEYWORD('which-jobs-supported'),
    ipp.KEYWORD('printer-get-attributes-supported'),
    ipp.KEYWORD('printer-uri-supported'),
    ipp.KEYWORD('media-supported'),
    ipp.KEYWORD('uri-security-supported'),
    ipp.KEYWORD('printer-dns-sd-name'),
    ipp.KEYWORD('printer-info'),
    ipp.KEYWORD('printer-name'),
    ipp.KEYWORD('printer-location'),
    ipp.KEYWORD('printer-make-and-model'),
    ipp.KEYWORD('color-supported'),
    ipp.KEYWORD('print-color-mode-supported'),
    ipp.KEYWORD('uri-authentication-supported'),
    ipp.KEYWORD('media-default'),
    ipp.KEYWORD('copies-default'),
    ipp.KEYWORD('print-color-mode-default'),
    ipp.KEYWORD('print-scaling-default'),
    ipp.KEYWORD('printer-uuid'),
    ipp.KEYWORD('printer-device-id'),
    ipp.KEYWORD('printer-geo-location'),
    ipp.KEYWORD('printer-up-time'),
    ipp.KEYWORD('queued-job-count'),
    ipp.KEYWORD('printer-is-accepting-jobs'),
    ipp.KEYWORD('printer-state'),
    ipp.KEYWORD('printer-state-reasons'),
    ipp.KEYWORD('printer-more-info'),
    ipp.KEYWORD('printer-icons'),
    ipp.KEYWORD('media-col-default'),
    ipp.KEYWORD('media-size-supported'),
    ipp.KEYWORD('printer-kind'),
    ipp.KEYWORD('multiple-operation-time-out-action'),
    ipp.KEYWORD('printer-organization'),
    ipp.KEYWORD('printer-organizational-unit'),
    ipp.KEYWORD('printer-state-message'),
    ipp.KEYWORD('printer-alert'),
    ipp.KEYWORD('printer-alert-description'),
    ipp.KEYWORD('document-format-version-supplied'),
    ipp.KEYWORD('printer-config-change-date-time'),
    ipp.KEYWORD('printer-config-change-time'),
    ipp.KEYWORD('printer-state-change-date-time'),
    ipp.KEYWORD('printer-state-change-time'),
    ipp.KEYWORD('printer-current-time'),
    ipp.KEYWORD('printer-firmware-name'),
    ipp.KEYWORD('printer-firmware-string-version'),
    ipp.KEYWORD('printer-firmware-version'),
    ipp.KEYWORD('urf-supported'),
    ipp.KEYWORD('overrides-supported'),
    ipp.KEYWORD('job-mandatory-attributes'),
    ipp.KEYWORD('printer-fax-log-uri'),
    ipp.KEYWORD('feed-orientation-default'),
    ipp.KEYWORD('feed-orientation-supported'),
    ipp.KEYWORD('finishings-default'),
    ipp.KEYWORD('finishings-supported'),
    ipp.KEYWORD('jpeg-x-dimension-supported'),
    ipp.KEYWORD('jpeg-y-dimension-supported'),
    ipp.KEYWORD('jpeg-k-octets-supported'),
    ipp.KEYWORD('pdf-k-octets-supported'),
    ipp.KEYWORD('marker-colors'),
    ipp.KEYWORD('marker-high-levels'),
    ipp.KEYWORD('marker-levels'),
    ipp.KEYWORD('marker-low-levels'),
    ipp.KEYWORD('marker-names'),
    ipp.KEYWORD('marker-types'),
    ipp.KEYWORD('media-col-ready'),
    ipp.KEYWORD('media-ready'),
    ipp.KEYWORD('media-source-supported'),
    ipp.KEYWORD('media-type-supported'),
    ipp.KEYWORD('landscape-orientation-requested-preferred'),
    ipp.KEYWORD('output-bin-default'),
    ipp.KEYWORD('output-bin-supported'),
    ipp.KEYWORD('pages-per-minute'),
    ipp.KEYWORD('printer-input-tray'),
    ipp.KEYWORD('printer-output-tray'),
    ipp.KEYWORD('printer-supply-info-uri'),
    ipp.KEYWORD('pwg-raster-document-sheet-back'),
    ipp.KEYWORD('pwg-raster-document-resolution-supported'),
    ipp.KEYWORD('pwg-raster-document-type-supported'),
    ipp.KEYWORD('sides-default'),
    ipp.KEYWORD('sides-supported'),
    ipp.KEYWORD('job-pages-per-set-supported')
  ],
  'printer-config-change-date-time': ipp.DATE('2026-07-06T16:35:48+03:00'),
  'printer-config-change-time': ipp.INTEGER(1783344948),
  'printer-state-change-date-time': ipp.DATE('2026-07-06T16:35:48+03:00'),
  'printer-state-change-time': ipp.INTEGER(1783344948),
  'printer-current-time': ipp.DATE('2026-07-06T16:35:48+03:00'),
  'media-col-database': [
    {
      'media-size': {'x-dimension': ipp.INTEGER(21000), 'y-dimension': ipp.INTEGER(29700)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400)
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(14800), 'y-dimension': ipp.INTEGER(21000)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(10500), 'y-dimension': ipp.INTEGER(14800)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(17600), 'y-dimension': ipp.INTEGER(25000)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(21590), 'y-dimension': ipp.INTEGER(35560)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400)
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(21590), 'y-dimension': ipp.INTEGER(27940)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400)
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(18415), 'y-dimension': ipp.INTEGER(26670)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(13970), 'y-dimension': ipp.INTEGER(21590)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400)
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(16200), 'y-dimension': ipp.INTEGER(22900)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(11400), 'y-dimension': ipp.INTEGER(16200)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(11000), 'y-dimension': ipp.INTEGER(22000)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(9842), 'y-dimension': ipp.INTEGER(19050)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(18200), 'y-dimension': ipp.INTEGER(25700)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400)
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(12800), 'y-dimension': ipp.INTEGER(18200)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(10500), 'y-dimension': ipp.INTEGER(23500)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(10000), 'y-dimension': ipp.INTEGER(14800)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(14800), 'y-dimension': ipp.INTEGER(20000)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(19685), 'y-dimension': ipp.INTEGER(27305)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400)
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(21590), 'y-dimension': ipp.INTEGER(33020)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(10477), 'y-dimension': ipp.INTEGER(24130)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(9842), 'y-dimension': ipp.INTEGER(22542)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(9207), 'y-dimension': ipp.INTEGER(16510)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400),
      'media-source': ipp.KEYWORD('by-pass-tray')
    },
    {
      'media-size': {'x-dimension': ipp.INTEGER(21000), 'y-dimension': ipp.INTEGER(33000)},
      'media-left-margin': ipp.INTEGER(400),
      'media-right-margin': ipp.INTEGER(400),
      'media-top-margin': ipp.INTEGER(400),
      'media-bottom-margin': ipp.INTEGER(400)
    }
  ]
}


# eSCL scanner parameters:
escl.caps = {
  'Version': '2.62',
  'MakeAndModel': 'Kyocera ECOSYS M2040dn',
  'SerialNumber': 'VCF9192281',
  'Uuid': UUID('4509a320-00a0-008f-00b6-002507510eca'),
  'AdminUri': 'https://KM7B6A91.local/airprint',
  'IconUri': 'https://KM7B6A91.local/printer-icon/machine_128.png',
  'Platen': {
    'PlatenInputCaps': {
      'MinWidth': 118,
      'MaxWidth': 2551,
      'MinHeight': 118,
      'MaxHeight': 3508,
      'SupportedIntents': ['Document', 'TextAndGraphic', 'Photo', 'Preview'],
      'SettingProfiles': [
        {
          'ColorModes': ['BlackAndWhite1', 'Grayscale8', 'RGB24'],
          'DocumentFormats': ['image/jpeg', 'application/pdf'],
          'SupportedResolutions': [
            {
              'DiscreteResolutions': [
                {'XResolution': 200, 'YResolution': 100},
                {'XResolution': 200, 'YResolution': 200},
                {'XResolution': 200, 'YResolution': 400},
                {'XResolution': 300, 'YResolution': 300},
                {'XResolution': 400, 'YResolution': 400},
                {'XResolution': 600, 'YResolution': 600}
              ]
            }
          ]
        }
      ],
      'FeedDirections': ['ShortEdgeFeed', 'LongEdgeFeed']
    }
  },
  'Adf': {
    'AdfSimplexInputCaps': {
      'MinWidth': 591,
      'MaxWidth': 2551,
      'MinHeight': 591,
      'MaxHeight': 4205,
      'SupportedIntents': ['Document', 'TextAndGraphic', 'Photo', 'Preview'],
      'SettingProfiles': [
        {
          'ColorModes': ['BlackAndWhite1', 'Grayscale8', 'RGB24'],
          'DocumentFormats': ['image/jpeg', 'application/pdf'],
          'SupportedResolutions': [
            {
              'DiscreteResolutions': [
                {'XResolution': 200, 'YResolution': 100},
                {'XResolution': 200, 'YResolution': 200},
                {'XResolution': 200, 'YResolution': 400},
                {'XResolution': 300, 'YResolution': 300},
                {'XResolution': 400, 'YResolution': 400},
                {'XResolution': 600, 'YResolution': 600}
              ]
            }
          ]
        }
      ]
    },
    'AdfDuplexInputCaps': {
      'MinWidth': 591,
      'MaxWidth': 2551,
      'MinHeight': 591,
      'MaxHeight': 4205,
      'SupportedIntents': ['Document', 'TextAndGraphic', 'Photo', 'Preview'],
      'SettingProfiles': [
        {
          'ColorModes': ['BlackAndWhite1', 'Grayscale8', 'RGB24'],
          'DocumentFormats': ['image/jpeg', 'application/pdf'],
          'SupportedResolutions': [
            {
              'DiscreteResolutions': [
                {'XResolution': 200, 'YResolution': 100},
                {'XResolution': 200, 'YResolution': 200},
                {'XResolution': 200, 'YResolution': 400},
                {'XResolution': 300, 'YResolution': 300},
                {'XResolution': 400, 'YResolution': 400},
                {'XResolution': 600, 'YResolution': 600}
              ]
            }
          ]
        }
      ],
      'FeedDirections': ['ShortEdgeFeed', 'LongEdgeFeed']
    },
    'FeederCapacity': 75,
    'AdfOptions': ['DetectPaperLoaded', 'SelectSinglePage', 'Duplex']
  },
  'CompressionFactorSupport': {
    'Min': 1,
    'Max': 5,
    'Normal': 1,
    'Step': 1
  },
  'SharpenSupport': {
    'Min': -3,
    'Max': 3,
    'Normal': 0,
    'Step': 1
  }
}


# ----- PUT YOUR ESCL HOOKS HERE -----

# Called on request:  POST /{root}/ScanJobs
#
# def escl_onScanJobsRequest (q: query.Query, rq: escl.ScanSettings):

# Called on response: GET /{JobUri}/NextDocument
#
# def escl_onNextDocumentResponse (q: query.Query, flt: escl.ImageFilter):

