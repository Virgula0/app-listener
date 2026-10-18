# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 468 files · ~1,648,204 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6136 nodes · 25921 edges · 224 communities (194 shown, 30 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2422 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e2b29e35`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- buildTrustedSet
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- mkdirs
- selection.go
- Network
- protectedPage
- github.com/charmbracelet/bubbletea.Cmd
- .launch
- os.FileMode
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- setpassword_authz_test.go
- globBuilder
- guard.bpf.c
- fscrypt/orphans.go
- watchSet
- icon
- headless.go
- UI
- serve.go
- IntegrationSuite
- highlight.go
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- backend
- controlServer
- fileedit_model_test.go
- os.File
- IntegrationSuite
- absPath
- Monitor
- NetEventType
- GuardEvent
- Open
- stdio
- editView
- fadeScroll
- LibraryClosure
- GenerateConf
- proto.go
- daemonconfig.go
- desktop.go
- TrustGuard
- backups.go
- update_test.go
- GuardInodeKey
- wipe
- Vault
- navItem
- exe_supersede.h
- audit.go
- install/install.go
- bpfstats/main.go
- gui/session.go
- RawTarget
- update/update.go
- changelog.go
- fcntl
- Vault
- Guard
- sync.Mutex
- btrfsLayoutFrom
- filevault_test.go
- tempGrantSetup
- RefreshCatalog
- runMonitor
- process_vm_readv.c
- CandidateDir
- filevault.go
- planUpdaters
- .decryptFileWithProgress
- github.com/charmbracelet/bubbles/textarea.Model
- runView
- Finder
- inodeChain
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- DaemonEvent
- .FilterExistingWhitelist
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- Candidate
- binaryVetter
- NetGuardEvent
- Well-known bypass classes
- NewDumpHook
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- readGrantRequest
- Build & sign release assets (reusable workflow)
- installPage
- golang.org/x/sys/unix.Stat_t
- mountsUnitTest
- NetworkMonitorRepository
- app-listener Terminal Demo Recording
- IntegrationSuite
- shQuote
- uninstall_system.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- expandPlaceholders
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- TestBatchInodeRows
- make test-integration (rootful Docker suite)
- protected_test.go
- supersede_probe.c
- showConfirm
- trace-app-libs.sh
- IntegrationSuite
- options.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- fakeGuardRepo
- time.Duration
- watchPattern
- runBackend
- catalogRefresher
- net_tester/main.go
- edit-protected
- uninstall/uninstall.go
- parseGuardEvent
- logGrid
- ConfinedEntry
- EventType
- layout.go
- githubRelease
- journal
- Installation
- backend_test.go
- graphify-refresh.sh
- vaultfs_test.go
- newTestTarget
- ParseEventsFlag
- pinstate.go
- feed.go
- shell
- DumpHook
- WithSuspectsFrom
- BinaryStat
- fileedit_symlink_test.go
- sanitizeTerminalText
- Limitations
- Troubleshooting
- IntegrationSuite
- FetchLatest
- Daemon configuration
- KernelDev
- systemPatterns
- systemd.go
- Compatibility
- How it works
- configedit_test.go
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- bpfLsmListed
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- netModel
- checkRunningBinaryMatchesInstalled
- ResolvePinBase
- monitor: see what touches a path
- kernel-prof.sh
- controlManager
- pathCache
- Execute
- Config
- writeWithin
- runNetworkGuard
- copyRegularAt
- Desktop GUI
- IntegrationSuite
- Ledger
- time.Time
- .eventTypeAllowed
- changelog_test.go
- TestEditTextMapsUICancel
- fakeDaemon
- buildOneGuard
- SuperblockDev
- loadDaemonConfig
- bleedLayout
- Vault
- TestMain
- launchKey
- parsePasswd
- inode_dev.h

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `GuardInodeKey` - 96 edges
3. `Config` - 95 edges
4. `Load()` - 89 edges
5. `shell` - 83 edges
6. `IntegrationSuite` - 82 edges
7. `absPath()` - 81 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go

## Import Cycles
- None detected.

## Communities (224 total, 30 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.12
Nodes (24): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey, Check() (+16 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (123): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod() (+115 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.12
Nodes (42): rcTarget, invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), BunTmpDir(), bunTargets(), createOwned(), EnsureAddKeysToAgent() (+34 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.04
Nodes (34): addErrKind, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+26 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "selection.go"
Cohesion: 0.13
Nodes (15): pos, cellToRune(), clipboardText(), cursorPos(), Selection, offset(), setClipboard(), cellOf() (+7 more)

### Community 11 - "Network"
Cohesion: 0.08
Nodes (15): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+7 more)

### Community 12 - "protectedPage"
Cohesion: 0.06
Nodes (21): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, sizedDialog, Event (+13 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (13): diffModel, editorModel, renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), RunFileEditorSession(), RunSession(), ParseMode() (+5 more)

### Community 14 - ".launch"
Cohesion: 0.31
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 15 - "os.FileMode"
Cohesion: 0.08
Nodes (33): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), dupFD() (+25 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (49): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+41 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (12): ownerPaths(), backingDeviceUnion(), Resource, BackingDevices(), GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase (+4 more)

### Community 23 - "setpassword_authz_test.go"
Cohesion: 0.12
Nodes (37): CheckPassword(), Hash(), HashConfigured(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+29 more)

### Community 24 - "globBuilder"
Cohesion: 0.15
Nodes (12): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, compileGlobNames(), globText() (+4 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "fscrypt/orphans.go"
Cohesion: 0.19
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 27 - "watchSet"
Cohesion: 0.30
Nodes (4): isSymlink(), watchDir, watchPos, watchSet

### Community 28 - "icon"
Cohesion: 0.08
Nodes (30): binaryList, filesPage, handButton, homePage, multiAnswer, newBinaryList(), newBinaryRow(), TestAdaptiveGridDropsColumns() (+22 more)

### Community 29 - "headless.go"
Cohesion: 0.29
Nodes (8): stagedBinary, stageRelease(), Check(), CheckGUI(), inspect(), buildTiny(), TestCheck(), TestCheckGUI()

### Community 30 - "UI"
Cohesion: 0.08
Nodes (22): showChangelog(), channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses() (+14 more)

### Community 31 - "serve.go"
Cohesion: 0.09
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 32 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 33 - "highlight.go"
Cohesion: 0.21
Nodes (6): fg(), Highlighter, lexerFor(), shebangLexer(), TestTokenizeKeepsLineCount(), tokenClassOf()

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.12
Nodes (24): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+16 more)

### Community 37 - "backend"
Cohesion: 0.14
Nodes (7): checkSecret(), backend, backendOp, backendState, Message, Request, RequireDaemonStopped()

### Community 38 - "controlServer"
Cohesion: 0.17
Nodes (8): readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), controlServer, editControlSession, grantRequest

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 40 - "os.File"
Cohesion: 0.04
Nodes (61): multicallAdmissible(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched() (+53 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 43 - "Monitor"
Cohesion: 0.08
Nodes (7): NewMonitor(), newPathCache(), dirTarget(), newMountTable(), Monitor, monitorUnitTest, watchedInode

### Community 44 - "NetEventType"
Cohesion: 0.06
Nodes (23): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), Cstr(), FormatAddr() (+15 more)

### Community 45 - "GuardEvent"
Cohesion: 0.06
Nodes (20): holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+12 more)

### Community 46 - "Open"
Cohesion: 0.18
Nodes (19): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+11 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "editView"
Cohesion: 0.12
Nodes (7): codeEntry, editView, Entry, selectable(), fileSize(), joinRel(), relOf()

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+19 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.12
Nodes (16): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+8 more)

### Community 51 - "GenerateConf"
Cohesion: 0.06
Nodes (54): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, Section, findSectionEnd(), findSectionStart(), firstContentLine() (+46 more)

### Community 52 - "proto.go"
Cohesion: 0.07
Nodes (20): Client, File, User, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestCleanRel() (+12 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.13
Nodes (39): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+31 more)

### Community 54 - "desktop.go"
Cohesion: 0.10
Nodes (14): releaseBuild(), Run(), encryptDirectories(), TestProgressReader(), stageItem, stageState, textAnswer, EnsureSystemSetup() (+6 more)

### Community 55 - "TrustGuard"
Cohesion: 0.07
Nodes (17): tempGrant, tempMaskOp, tempRow, trustHook, trustRows, deleteInoKeys(), deleteResKeys(), setSupersedeTrust() (+9 more)

### Community 56 - "backups.go"
Cohesion: 0.18
Nodes (20): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), Delete(), Find(), find(), Backup (+12 more)

### Community 57 - "update_test.go"
Cohesion: 0.09
Nodes (16): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestAssetsFor(), TestFetchGUIWithoutAsset(), TestGUIAssetsFor(), testKeyPair() (+8 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (57): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), ensureBtrfsLayout(), nextLabel() (+49 more)

### Community 59 - "wipe"
Cohesion: 0.12
Nodes (17): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn() (+9 more)

### Community 60 - "Vault"
Cohesion: 0.08
Nodes (40): CheckBPFLSM(), buildConcurrency(), buildGuards(), configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), makeReloadHandler(), markReady() (+32 more)

### Community 61 - "navItem"
Cohesion: 0.03
Nodes (28): appTheme, binaryRow, buttonCursor, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer (+20 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 63 - "audit.go"
Cohesion: 0.25
Nodes (14): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+6 more)

### Community 64 - "install/install.go"
Cohesion: 0.35
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 65 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 66 - "gui/session.go"
Cohesion: 0.07
Nodes (41): ancestorEnv(), cleanupSessionEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField() (+33 more)

### Community 67 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+2 more)

### Community 68 - "update/update.go"
Cohesion: 0.18
Nodes (24): applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), confirmUnit(), deploy(), deployGUI(), downloadAndVerify() (+16 more)

### Community 69 - "changelog.go"
Cohesion: 0.12
Nodes (39): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+31 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "Vault"
Cohesion: 0.23
Nodes (3): Vault, Open(), Unopened()

### Community 72 - "Guard"
Cohesion: 0.05
Nodes (31): deferredBinary, exeRow, rootHandle, snapNet, BinaryEntry, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys() (+23 more)

### Community 73 - "sync.Mutex"
Cohesion: 0.09
Nodes (22): Target, guardPill(), newActivity(), newExeResolver(), FileEngines, Files, FileSpec, FileTarget (+14 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "filevault_test.go"
Cohesion: 0.15
Nodes (27): harnessSuite, looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated() (+19 more)

### Community 76 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 77 - "RefreshCatalog"
Cohesion: 0.16
Nodes (15): sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), refreshCatalogConfig(), RefreshCatalog() (+7 more)

### Community 78 - "runMonitor"
Cohesion: 0.11
Nodes (15): CheckEBPF(), runHeadless(), runMonitor(), runTUI(), runNetworkMonitor(), runTUI(), PrintLogo(), NewModel() (+7 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "CandidateDir"
Cohesion: 0.14
Nodes (17): cleanOrphanedMetadata(), Discover(), DiscoverForUsers(), DiscoverSystem(), CandidateDir, TestDiscoverForUsersDeduplicates(), TestDiscoverOnlyExisting(), TestDiscoverPerUserHomes() (+9 more)

### Community 81 - "filevault.go"
Cohesion: 0.27
Nodes (13): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+5 more)

### Community 82 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 83 - ".decryptFileWithProgress"
Cohesion: 0.16
Nodes (8): entryOf(), migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow(), FileModeToUnixMode()

### Community 84 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (12): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+4 more)

### Community 86 - "Finder"
Cohesion: 0.13
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 87 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 88 - "liveSession"
Cohesion: 0.19
Nodes (11): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents() (+3 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "DaemonEvent"
Cohesion: 0.08
Nodes (39): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+31 more)

### Community 92 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (11): foreignPinnedLink(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+3 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "Candidate"
Cohesion: 0.25
Nodes (14): appendSectionsAndEdit(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates() (+6 more)

### Community 100 - "binaryVetter"
Cohesion: 0.11
Nodes (14): newBinaryVetter(), openBinaryVetter(), newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), startTrustGuard(), trustStartupError(), binaryVetter (+6 more)

### Community 101 - "NetGuardEvent"
Cohesion: 0.16
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "readGrantRequest"
Cohesion: 0.18
Nodes (13): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+5 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.08
Nodes (14): binSource, installedBinary, installPage, latestRelease, uninstallCard, updatePage, hashFile(), hashReader() (+6 more)

### Community 110 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 111 - "mountsUnitTest"
Cohesion: 0.12
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+2 more)

### Community 112 - "NetworkMonitorRepository"
Cohesion: 0.11
Nodes (15): NetworkMonitorRepository, NewGuardUseCase(), NewMonitorUseCase(), NetworkMonitorUseCase, NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo() (+7 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "shQuote"
Cohesion: 0.15
Nodes (4): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite

### Community 116 - "uninstall_system.go"
Cohesion: 0.26
Nodes (12): detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertBunLaunchers(), revertSSHAgents() (+4 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "expandPlaceholders"
Cohesion: 0.15
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "TestBatchInodeRows"
Cohesion: 0.40
Nodes (6): inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), InodeOccupancy(), TestCountInodeRows()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "protected_test.go"
Cohesion: 0.26
Nodes (6): fakeLink, isOwnExecutable(), newFakeLink(), protectedShell(), TestForwardRefusesOwnExecutable(), TestProtectedEditFlow()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "showConfirm"
Cohesion: 0.20
Nodes (9): handCursors(), handCursorsOn(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog(), showForm(), shimmedButtons() (+1 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "options.go"
Cohesion: 0.24
Nodes (10): installAction, InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks (+2 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "fakeGuardRepo"
Cohesion: 0.05
Nodes (24): confirmUserPlaced(), forwardCheck(), SupersededBinary, binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps() (+16 more)

### Community 143 - "time.Duration"
Cohesion: 0.17
Nodes (12): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+4 more)

### Community 144 - "watchPattern"
Cohesion: 0.24
Nodes (10): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern, Process inspectors, rootOwnedInode() (+2 more)

### Community 145 - "runBackend"
Cohesion: 0.13
Nodes (17): reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), runEditProtected(), chooseResource(), promptPassword(), promptPasswordWith() (+9 more)

### Community 146 - "catalogRefresher"
Cohesion: 0.15
Nodes (8): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), catalogRefresher, debouncer

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 149 - "uninstall/uninstall.go"
Cohesion: 0.05
Nodes (53): confirmMasterKeyOverwrite(), runGenKey(), TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), ensureMasterKey(), prepareInstallation(), collectFilesystemPrereqs() (+45 more)

### Community 150 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.13
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 154 - "EventType"
Cohesion: 0.08
Nodes (12): snapEngines, snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+4 more)

### Community 155 - "layout.go"
Cohesion: 0.29
Nodes (7): adaptiveGrid, checkGrid, gapVBox, cellMin(), newGapVBox(), rowHeights(), visible()

### Community 156 - "githubRelease"
Cohesion: 0.15
Nodes (24): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), tags(), TestReleasesSince(), compareStableVersions(), filterChannel() (+16 more)

### Community 157 - "journal"
Cohesion: 0.14
Nodes (7): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.28
Nodes (18): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+10 more)

### Community 161 - "vaultfs_test.go"
Cohesion: 0.14
Nodes (20): DefaultNewFileMode(), IsBinaryContent(), IsOctalMode(), inodeOf(), mustRead(), mustWrite(), openRootT(), TestDefaultNewFileMode() (+12 more)

### Community 162 - "newTestTarget"
Cohesion: 0.33
Nodes (11): newTestTarget(), since(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst(), TestApplyOffStopsGuard(), TestApplyOffToWhitelistStartsGuard() (+3 more)

### Community 163 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 164 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (47): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+39 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (26): pageDef, pageID, resizeWatch, runHost, shell, sysStatus, newLogo(), newShell() (+18 more)

### Community 167 - "DumpHook"
Cohesion: 0.16
Nodes (8): logHook, logLine, logRow(), logStyle(), formatEntry(), sortedKeys(), sortStrings(), DumpHook

### Community 168 - "WithSuspectsFrom"
Cohesion: 0.43
Nodes (3): WithSuspectsFrom(), netguardHarness, Option

### Community 169 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 170 - "fileedit_symlink_test.go"
Cohesion: 0.23
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 171 - "sanitizeTerminalText"
Cohesion: 0.10
Nodes (19): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+11 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "FetchLatest"
Cohesion: 0.15
Nodes (16): installHooks(), staged(), downloadFile(), downloadReleaseFiles(), FetchLatest(), fetchRelease(), fetchReleases(), FetchVerified() (+8 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "KernelDev"
Cohesion: 0.10
Nodes (24): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+16 more)

### Community 178 - "systemPatterns"
Cohesion: 0.33
Nodes (6): systemPatterns(), systemRule(), parseEvents(), specFromForm(), ParseEventType(), BinaryRule

### Community 179 - "systemd.go"
Cohesion: 0.06
Nodes (62): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), applyLiveRefresh(), buildBinaryIfNeeded() (+54 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "bpfLsmListed"
Cohesion: 0.40
Nodes (6): bpfLsmListed(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "netModel"
Cohesion: 0.18
Nodes (9): syncViewport(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg (+1 more)

### Community 195 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 196 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "controlManager"
Cohesion: 0.22
Nodes (4): newControlManager(), startControlManager(), configEditor, controlManager

### Community 200 - "pathCache"
Cohesion: 0.60
Nodes (3): evalSymlinksOrEmpty(), pathCache, pathCacheEntry

### Community 201 - "Execute"
Cohesion: 0.29
Nodes (7): Requested(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger(), subcommand(), main()

### Community 202 - "Config"
Cohesion: 0.09
Nodes (40): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions() (+32 more)

### Community 203 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 204 - "runNetworkGuard"
Cohesion: 0.12
Nodes (16): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI() (+8 more)

### Community 205 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "Ledger"
Cohesion: 0.10
Nodes (12): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, applySchema() (+4 more)

### Community 210 - "time.Time"
Cohesion: 0.24
Nodes (5): LockedError, newWatchSet(), fakeInotify, sweepID, sweepMemo

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "changelog_test.go"
Cohesion: 0.31
Nodes (8): newChangelogModel(), rel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog(), TestNewChangelogModelViewContainsTitleAndNotes()

### Community 216 - "buildOneGuard"
Cohesion: 0.07
Nodes (43): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), displayEntries() (+35 more)

### Community 217 - "SuperblockDev"
Cohesion: 0.29
Nodes (7): FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev(), SuperblockDevPath(), TestMountinfoDev(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 219 - "loadDaemonConfig"
Cohesion: 0.29
Nodes (7): loadDaemonConfig(), noteReloadedSettings(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 222 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 227 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 232 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 568 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **30 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `edit-protected` to `backups.go`, `backend`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11739965625315944 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02311111111111111 - nodes in this community are weakly interconnected._
- **Why does `Delete()` connect `backups.go` to `edit-protected`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08957708957708958 - nodes in this community are weakly interconnected._