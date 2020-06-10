import { Commit } from './blamelist-panel-sk';

export const fakeNow = Date.parse('2020-03-22T00:00:00.000Z');

export const blamelist19: Commit[] = [
  {
    commit_time: 1584835000,
    hash: 'dded3c7506efc5635e60ffb7a908cbe8f1f028f1',
    author: 'Alfa (alfa@example.com)',
    message: 'Update provisioning_profile to unbreak iOS since cert refresh',
    is_cl: false,
  },
  {
    commit_time: 1584834200,
    hash: '44fc53b7f5d8390d7ae228d5fa8bac0b225ed2b9',
    author: 'Bravo (bravo@example.com)',
    message: 'Add BGR_10A2 support to Ganesh',
    is_cl: false,
  },
  {
    commit_time: 1584834100,
    hash: 'c50f4e9399cc7ae72b05eb641b50a7cd06d5af82',
    author: 'Charlie (charlie@example.com)',
    message: 'Simplify GrTextBlob::flush',
    is_cl: false,
  },
  {
    commit_time: 1584834000,
    hash: '9145f784f3261f227846e5b08dc2691a888b113c',
    author: 'Delta (deltadelta@example.com)',
    message: 'Implement D3D copySurface.',
    is_cl: false,
  },
  {
    commit_time: 1584824000,
    hash: '02b91385f58b84dc14cb6f81478b36d7fcb0d24b',
    author: 'Echo Foxtrot (ef@example.com)',
    message: 'Rename flush -> addOp',
    is_cl: false,
  },
  {
    commit_time: 1584814000,
    hash: 'a1e90aba441b0e68a07e3ab6660cd673335a76f5',
    author: 'Golf Golf Golf Golf Golf Golf Golf(golf@example.com)',
    message: 'Turn on Vulkan bots for Galaxy S20.',
    is_cl: false,
  },
  {
    commit_time: 1584804000,
    hash: 'b2b9fbc5675c46cb3af30ccd9b6e20bc57ec04c8',
    author: 'Hotel India (hotel.ind@example.com)',
    message: 'Compare all fields in SkSL::Layout::operator==',
    is_cl: false,
  },
  {
    commit_time: 1584800000,
    hash: 'f1792cde0b2b7b612cac6aa5b9e7a8ce16b9b2cb',
    author: 'Hotel India (hotel.ind@example.com)',
    message: 'Use constant swizzle syntax in GrDrawVerticesOp',
    is_cl: false,
  },
  {
    commit_time: 1584780000,
    hash: 'fc75d5c52da617b06c7566c9e389587edde3f284',
    author: 'Juliett Kilo (julia.kilo@example.com)',
    message: '[skottie] Brightness and Contrast effect',
    is_cl: false,
  },
  {
    commit_time: 1584760000,
    hash: 'e15088801e9722d9a320463554a3a21839213b48',
    author: 'Mike Lima (lima@example.com)',
    message: 'detect failed matrix update in SkDraw::drawAtlas()',
    is_cl: false,
  },
  {
    commit_time: 1584760000,
    hash: '548de7451e752ce0d6fd842c1bb4f04af4c0afdc',
    author: 'Hotel India (hotel.ind@example.com)',
    message: 'Change Marker IDs to be strings',
    is_cl: false,
  },
  {
    commit_time: 1584740000,
    hash: '7676f4ecf3ec2f544701465faca1d33af1489822',
    author: 'Zulu (zulu@example.com)',
    message: 'Remove SkFrontBufferedStream',
    is_cl: false,
  },
  {
    commit_time: 1584720000,
    hash: '4baa7326ccfbfaf66e25d91628378536d4999999',
    author: 'Papa November (pnov@example.com)',
    message: '[infra] Add POC task driver',
    is_cl: false,
  },
  {
    commit_time: 1584700000,
    hash: 'f5132a05c893a86f8bf26bcf3253985d9973fea2',
    author: 'Quebec (quebec@example.com)',
    message: "Reland \"Optimize GrTessellatePathOp's code to emit inner tri",
    is_cl: false,
  },
  {
    commit_time: 1584600005,
    hash: '3d599d16ecd8009dfb185522e87c5c1cf8a47057',
    author: 'skia-recreate-skps (skia-recreate-skps@example.com)',
    message: 'Update Go Deps',
    is_cl: false,
  },
  {
    commit_time: 1584600005,
    hash: '7a0517f8bdc0d168565a2944f0c14db06084384d',
    author: 'skia-autoroll (skia-autoroll@example.com)',
    message: 'Roll third_party/externals/angle2 3cb9c4bee9b3..4395170e6091',
    is_cl: false,
  },
  {
    commit_time: 1584600003,
    hash: 'd18100afb62401b49e0b3a5c87ed8fa006e8d3a6',
    author: 'skia-autoroll (skia-autoroll@example.com)',
    message: 'Roll ../src 9781ff27c9e9..78824aa9d99f (468 commits)',
    is_cl: false,
  },
  {
    commit_time: 1584600002,
    hash: '2ac6fa0ae796b7f050bd8e8ad60890427a4ec15e',
    author: 'skia-autoroll (skia-autoroll@example.com)',
    message: 'Roll third_party/externals/swiftshader 60aa34a990fa..2717702',
    is_cl: false,
  },
  {
    commit_time: 1584600001,
    hash: '667edf14ad72966ec36aa6cd705b98cb7d7eee28',
    author: 'skia-autoroll (skia-autoroll@example.com)',
    message: 'Roll third_party/externals/dawn 00b90ea83262..88f2ec853f80 (',
    is_cl: false,
  },
];

export const clBlamelist: Commit[] = [
  {
    commit_time: 1584600001,
    hash: '12345',
    author: 'skia-autoroll (skia-autoroll@example.com)',
    message: 'Roll third_party/externals/dawn 00b90ea83262..88f2ec853f80 (',
    is_cl: true,
  },
];