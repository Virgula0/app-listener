# Graph Report - app-listener  (2026-10-28)

## Corpus Check
- 526 files · ~1,742,854 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6839 nodes · 30081 edges · 246 communities (201 shown, 45 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2776 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `525f0089`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- BinaryStat
- New
- protectedPage
- find_test.go
- go_pkg_path_filepath
- installPage
- IntegrationSuite
- guard_trust.bpf.c
- TrustGuard
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- Guard
- guard.bpf.c
- github.com/cilium/ebpf.Map
- github.com/charmbracelet/bubbles/textarea.Model
- FileTarget
- .step
- links_test.go
- TestMonitorUseCaseLifecycle
- BinaryEntry
- codeedit.go
- picker
- RefreshSection
- IntegrationSuite
- Message
- sync.Mutex
- safeio.go
- canonicalBinaryPath
- IntegrationSuite
- absPath
- monitorUnitTest
- runNetworkGuard
- fileedit_model_test.go
- Classify
- stdio
- os.File
- fadeScroll
- vaultfs.go
- fyne.io/fyne/v2/driver/desktop.Cursor
- Client
- watchGroup
- fyne.io/fyne/v2/widget.Label
- .apply
- KernelDev
- mutedText
- engine
- DurableStat
- os.FileInfo
- Recorder
- exe_supersede.h
- NetEventType
- pickfs.go
- updatePage
- wm.go
- IntegrationSuite
- .FilterExistingWhitelist
- liveSession
- editView
- runNetworkMonitor
- CandidateDir
- fcntl
- bpfstats/main.go
- Config
- codecheck/main.go
- rcUser
- update/update.go
- process_vm_readv.c
- GuardEvent
- op
- fscrypt/orphans.go
- maintainerKeeper
- net_tester/main.go
- render.go
- finishBinary
- planUpdaters
- runBackend
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- EventType
- Monitor
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- ConfinedEntry
- pinstate.go
- Well-known bypass classes
- NetGuardEvent
- jit_provenance.c
- inspector_probe.c
- recordops.c
- root.go
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- Vault
- mountTable
- runOfflineEdit
- recordFlow
- IntegrationSuite
- UnifiedDiff
- editauth_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- GenerateConf
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- install/install.go
- LibraryClosure
- trace-app-libs.sh
- IntegrationSuite
- resolveFilesystemPrereqs
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- catalogRefresher
- Serve
- gui/session.go
- listFacultatives
- Ledger
- Candidate
- github.com/spf13/cobra.Command
- plan_test.go
- buildTrustedSet
- fakeStore
- inodeChain
- IntegrationSuite
- buildOneGuard
- procstats.go
- fileedit_symlink_test.go
- Installation
- backend_test.go
- graphify-refresh.sh
- layout.go
- IntegrationSuite
- .populateWatchInodes
- feed.go
- shell
- newTestTarget
- time.Time
- io.Writer
- Daemon configuration
- .showSized
- Limitations
- watchPattern
- IntegrationSuite
- runGuard
- netinject.c
- CopyTreeWithProgress
- parseGuardEvent
- TestMain
- Compatibility
- How it works
- ptrace_seize.c
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- logGrid
- RawTarget
- edit-protected
- promptPassword
- kernel-prof.sh
- sources_test.go
- pathCache
- run.go
- DaemonEvent
- Existing
- changelog_test.go
- find
- Desktop GUI
- IntegrationSuite
- github.com/charmbracelet/bubbletea.Cmd
- appTheme
- .eventTypeAllowed
- bleedLayout
- NewEventFanout
- parsePasswd
- .awaitDaemonUpOrDead
- configedit_test.go
- linkView
- tempGrantSetup
- SanitizeText
- shQuote
- github.com/stretchr/testify/suite.Suite
- prctl_gate
- compileGlobNames
- systemd.go
- supersede_probe.c
- sources
- testLimiter
- loadDaemonConfig
- logLine
- github.com/Virgula0/app-listener/tools/codecheck
- .exeKeys
- procRec
- taint_victim.c
- parseResize
- inode_dev.h
- Troubleshooting
- GeneralToolsOf
- removeKeyAndEmptyDir
- TestCleanRel
- github.com/cilium/ebpf.Program
- monitor: see what touches a path

## God Nodes (most connected - your core abstractions)
1. `Guard` - 123 edges
2. `Config` - 118 edges
3. `shell` - 91 edges
4. `Load()` - 90 edges
5. `absPath()` - 88 edges
6. `IntegrationSuite` - 83 edges
7. `fileEditModel` - 77 edges
8. `writeConfig()` - 69 edges
9. `IntegrationSuite` - 67 edges
10. `engine` - 61 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `vaultOpUnderGuard()` --calls--> `op`  [INFERRED]
  cmd/functions/install/steps.go → internal/maintainer/plan.go

## Import Cycles
- None detected.

## Communities (246 total, 45 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.12
Nodes (28): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+20 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (139): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestEveryInstallModeMapped(), TestForwardFlagSurface(), atomicWriteAt() (+131 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.12
Nodes (21): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr() (+13 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (77): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+69 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (45): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_user(), fill_identity() (+37 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (24): addErrKind, ReplacementCheck, classifyAddErr(), addRecorder(), copySelf(), guardUnitTest, runTool(), walkEntries() (+16 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "BinaryStat"
Cohesion: 0.23
Nodes (11): BinaryStat, binaryStatOf(), binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches() (+3 more)

### Community 11 - "New"
Cohesion: 0.08
Nodes (46): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+38 more)

### Community 12 - "protectedPage"
Cohesion: 0.06
Nodes (21): backendLink, checkGrid, configView, forwardForm, forwardView, protectedPage, recordForm, recordReview (+13 more)

### Community 13 - "find_test.go"
Cohesion: 0.16
Nodes (22): findAll(), lowerRunes(), NewFinder(), runesEqual(), findEditor(), keys(), TestFindAllSmartCase(), TestFinderLeavesSaveToHost() (+14 more)

### Community 14 - "go_pkg_path_filepath"
Cohesion: 0.36
Nodes (8): nextConfig(), recordTick, isAbsoluteCandidate(), TestCatalogSanity(), EditText(), LiveEmptyWhitelistRejected(), TestLiveEmptyWhitelistRejected(), hint

### Community 15 - "installPage"
Cohesion: 0.09
Nodes (9): actionGate, binSource, handCheck, installAction, installPage, uninstallCard, UninstallRequest, daemonStopNote() (+1 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.10
Nodes (56): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+48 more)

### Community 18 - "TrustGuard"
Cohesion: 0.10
Nodes (10): tempGrant, tempMaskOp, tempRow, trustRows, setSupersedeTrust(), liveTrust(), planTempBlock(), cStr() (+2 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): newFindInput(), clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.16
Nodes (46): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+38 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (14): ownerPaths(), fakeDaemon, Resource, GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase (+6 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "Guard"
Cohesion: 0.04
Nodes (32): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, deferredBinary, exeRow, Mode (+24 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (73): exe_refused(), inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade() (+65 more)

### Community 26 - "github.com/cilium/ebpf.Map"
Cohesion: 0.05
Nodes (37): TrustedDir, deleteInoKeys(), deleteResKeys(), supersedeUnnamed(), droppableFreed(), dropSupersededLocked(), exeGone(), freedSuperseded() (+29 more)

### Community 27 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.11
Nodes (18): findButton, pos, span, diffModel, column(), dedent(), lines(), moveTo() (+10 more)

### Community 28 - "FileTarget"
Cohesion: 0.11
Nodes (14): guardPill(), newExeResolver(), Files, FileSpec, FileTarget, Policy, guardRow(), newFileTarget() (+6 more)

### Community 29 - ".step"
Cohesion: 0.18
Nodes (13): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+5 more)

### Community 30 - "links_test.go"
Cohesion: 0.73
Nodes (5): checkLink(), linkFixture(), TestCheckLinksAccepts(), TestCheckLinksRefusesUserFiles(), TestCheckLinksSeesThePlannedTree()

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.17
Nodes (13): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+5 more)

### Community 32 - "BinaryEntry"
Cohesion: 0.06
Nodes (25): binaryList, netGuardTab, policyEditor, policyHooks, policyTab, BinaryEntry, diffBinaries(), TestDiffBinaries() (+17 more)

### Community 33 - "codeedit.go"
Cohesion: 0.11
Nodes (35): class, classOf(), Edit(), leadingIndent(), New(), onlySpaces(), SetText(), softTabBack() (+27 more)

### Community 34 - "picker"
Cohesion: 0.13
Nodes (6): keyEntry, picker, errText(), isFindShortcut(), newKeyEntry(), newPickRow()

### Community 35 - "RefreshSection"
Cohesion: 0.12
Nodes (23): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+15 more)

### Community 37 - "Message"
Cohesion: 0.11
Nodes (8): checkSecret(), File, backend, backendOp, backendState, Message, PwState, Request

### Community 38 - "sync.Mutex"
Cohesion: 0.08
Nodes (25): readConfigPut(), TestReadGrantRequestConfig(), newControlManager(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+17 more)

### Community 39 - "safeio.go"
Cohesion: 0.14
Nodes (22): installConfig(), installFileAs(), installUnit(), readAtNoFollow(), upsertFile(), upsertUnitAt(), writeAtNoFollow(), invokingIDs() (+14 more)

### Community 40 - "canonicalBinaryPath"
Cohesion: 0.14
Nodes (10): keptEntries, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath(), FirstOnBtrfs(), mountinfoDev(), OnBtrfs(), SuperblockDev() (+2 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.08
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 42 - "absPath"
Cohesion: 0.07
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "runNetworkGuard"
Cohesion: 0.15
Nodes (16): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+8 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "Classify"
Cohesion: 0.12
Nodes (16): Classify(), ValidateTargets(), fixtureExe(), mustParse(), privateTempRoots(), record(), TestClassifyLibs(), TestClassifyRefusesWorldWritableLocations() (+8 more)

### Community 48 - "os.File"
Cohesion: 0.04
Nodes (77): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+69 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (27): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+19 more)

### Community 50 - "vaultfs.go"
Cohesion: 0.12
Nodes (23): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead() (+15 more)

### Community 51 - "fyne.io/fyne/v2/driver/desktop.Cursor"
Cohesion: 0.10
Nodes (14): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), newButtonCursor(), newCheckGroup(), newRadioGroup() (+6 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (12): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+4 more)

### Community 53 - "watchGroup"
Cohesion: 0.08
Nodes (46): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+38 more)

### Community 54 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (21): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+13 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (15): Apply(), DeleteBackup(), Outcome, ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+7 more)

### Community 56 - "KernelDev"
Cohesion: 0.15
Nodes (17): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), codeVouched() (+9 more)

### Community 57 - "mutedText"
Cohesion: 0.08
Nodes (37): filesPage, homePage, multiAnswer, newLogo(), newBinaryList(), newBinaryRow(), newCheck(), newFadeVScroll() (+29 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (42): engine, inodeBatchFn, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates() (+34 more)

### Community 59 - "DurableStat"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 60 - "os.FileInfo"
Cohesion: 0.19
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 61 - "Recorder"
Cohesion: 0.14
Nodes (18): newBinary(), drain(), Recorder, Target, monitorTargets(), newRecorder(), opNames(), selfExe() (+10 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.16
Nodes (18): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+10 more)

### Community 63 - "NetEventType"
Cohesion: 0.06
Nodes (20): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+12 more)

### Community 64 - "pickfs.go"
Cohesion: 0.22
Nodes (14): listOpts, pickEntry, treeSearch, lstatOf(), outsideHardlinks(), plural(), entryIsDir(), filterEntries() (+6 more)

### Community 65 - "updatePage"
Cohesion: 0.09
Nodes (10): daemonPage, daemonState, installedBinary, latestRelease, runHost, updatePage, readInstalled(), newUpdatePage() (+2 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 69 - "liveSession"
Cohesion: 0.15
Nodes (15): dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward() (+7 more)

### Community 70 - "editView"
Cohesion: 0.10
Nodes (10): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+2 more)

### Community 71 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (12): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString() (+4 more)

### Community 72 - "CandidateDir"
Cohesion: 0.08
Nodes (23): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+15 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "Config"
Cohesion: 0.08
Nodes (47): resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startMaintainers(), auditGroupCoverage(), findResource(), groupWatchPaths(), askBunTmpdirUsers(), bunEntryConfigured() (+39 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "rcUser"
Cohesion: 0.12
Nodes (41): TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(), rcTarget, BunTmpDir(), TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent() (+33 more)

### Community 78 - "update/update.go"
Cohesion: 0.04
Nodes (91): installHooks(), staged(), changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), applyUpdate(), archAssets() (+83 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.04
Nodes (28): confirmUserPlaced(), forwardCheck(), holdHeadless(), HeadlessLine(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), commMatchesGuardedBinary() (+20 more)

### Community 81 - "op"
Cohesion: 0.16
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 82 - "fscrypt/orphans.go"
Cohesion: 0.18
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 83 - "maintainerKeeper"
Cohesion: 0.15
Nodes (14): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), grantNewBinaries(), maintainedApp (+6 more)

### Community 84 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 85 - "render.go"
Cohesion: 0.10
Nodes (17): lineStyles, viewportState, fg(), TestDaemonConfTokens(), TestTokenizeKeepsLineCount(), tokenClassOf(), baseStyle(), gutterLabel() (+9 more)

### Community 86 - "finishBinary"
Cohesion: 0.17
Nodes (13): beyondEvents(), classifyLibs(), finishBinary(), innermostResource(), isShellHelper(), libCovered(), noteWhitelisted(), worldWritableDir() (+5 more)

### Community 87 - "planUpdaters"
Cohesion: 0.19
Nodes (14): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+6 more)

### Community 88 - "runBackend"
Cohesion: 0.31
Nodes (7): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), changePasswordLive()

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "EventType"
Cohesion: 0.07
Nodes (18): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), specFromForm() (+10 more)

### Community 92 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.18
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.16
Nodes (3): eventTypeKey(), hook, NetGuard

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.17
Nodes (12): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+4 more)

### Community 100 - "ConfinedEntry"
Cohesion: 0.12
Nodes (30): markReady(), notifySystemdReady(), HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared() (+22 more)

### Community 101 - "pinstate.go"
Cohesion: 0.16
Nodes (24): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+16 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NetGuardEvent"
Cohesion: 0.11
Nodes (9): snapNetGuard, reasonText(), NetGuardEvent, TestReasonText(), protoString(), NetworkGuardRepository, netGuardEventLine, netGuardEventMsg (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.22
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.45
Nodes (13): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+5 more)

### Community 107 - "root.go"
Cohesion: 0.07
Nodes (30): releaseBuild(), Run(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary() (+22 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "Vault"
Cohesion: 0.23
Nodes (3): Vault, Open(), WriteFileKeepMeta()

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "runOfflineEdit"
Cohesion: 0.09
Nodes (22): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditTree(), coveredByWatch(), relOrPath(), TestAuditTree() (+14 more)

### Community 113 - "recordFlow"
Cohesion: 0.11
Nodes (25): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), formatElapsed(), lines(), offerBackupDeletion() (+17 more)

### Community 115 - "UnifiedDiff"
Cohesion: 0.14
Nodes (17): TestDaemonUnitWithMetadataOutput(), refreshCatalogConfig(), confirmUnit(), TestUnifiedDiff(), TestUnifiedDiffIdentical(), UnifiedDiff(), DaemonUnitWithMetadataOutput(), ConfirmUnitFunc (+9 more)

### Community 116 - "editauth_test.go"
Cohesion: 0.11
Nodes (25): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+17 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "GenerateConf"
Cohesion: 0.06
Nodes (55): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips() (+47 more)

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
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "install/install.go"
Cohesion: 0.06
Nodes (51): asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256(), prepareInstallation(), runInstall() (+43 more)

### Community 131 - "LibraryClosure"
Cohesion: 0.07
Nodes (28): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+20 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "resolveFilesystemPrereqs"
Cohesion: 0.20
Nodes (8): encryptDirectories(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), EnsureSystemSetup(), isKernelAtLeast54(), Prereq, fakeVault

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "catalogRefresher"
Cohesion: 0.09
Nodes (15): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), makeReloadHandler() (+7 more)

### Community 143 - "Serve"
Cohesion: 0.10
Nodes (26): basicAuth(), configureWebSocket(), isInteractiveTerminal(), newServeListener(), newServeServer(), readResizeLoop(), relayServeSignals(), runServe() (+18 more)

### Community 144 - "gui/session.go"
Cohesion: 0.21
Nodes (14): ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue(), TestOwnedByChecksTypeAndOwner() (+6 more)

### Community 145 - "listFacultatives"
Cohesion: 0.08
Nodes (37): installFacultatives(), installSSHAgentFacultative(), listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH() (+29 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 147 - "Candidate"
Cohesion: 0.17
Nodes (19): appendSectionsAndEdit(), addManualDirectories(), editConfig(), groupCandidates(), libraryBlocksFromCandidates(), pickDirectories(), pickFromCandidates(), runConfigEditor() (+11 more)

### Community 148 - "github.com/spf13/cobra.Command"
Cohesion: 0.09
Nodes (36): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe() (+28 more)

### Community 149 - "plan_test.go"
Cohesion: 0.27
Nodes (14): isWhiteout(), opaque(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles() (+6 more)

### Community 150 - "buildTrustedSet"
Cohesion: 0.18
Nodes (12): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+4 more)

### Community 152 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.08
Nodes (38): backingDeviceUnion(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon(), unlockPendingGroupRoots() (+30 more)

### Community 156 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 157 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.27
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (47): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+39 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (48): fakeLink, pageDef, pageID, reel, resizeWatch, shell, snapEngines, snapNet (+40 more)

### Community 167 - "newTestTarget"
Cohesion: 0.10
Nodes (24): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+16 more)

### Community 168 - "time.Time"
Cohesion: 0.13
Nodes (10): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID (+2 more)

### Community 169 - "io.Writer"
Cohesion: 0.13
Nodes (14): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), TestProgressReader(), fyneLogFilter (+6 more)

### Community 170 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 171 - ".showSized"
Cohesion: 0.04
Nodes (39): showChangelog(), InstallMode, InstallRequest, installUI, Options, sizedDialog, UninstallHooks, UpdateHooks (+31 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "watchPattern"
Cohesion: 0.20
Nodes (13): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), systemBinaryPatterns(), systemPatterns() (+5 more)

### Community 175 - "runGuard"
Cohesion: 0.08
Nodes (31): DrainUntilSignal(), CheckBPFLSM(), CheckEBPF(), MakeDisplayPaths(), configureDaemonLogging(), confirmMasterKeyOverwrite(), newPinGeneration(), prepareDaemonStart() (+23 more)

### Community 176 - "netinject.c"
Cohesion: 0.50
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 177 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 178 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 179 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "ptrace_seize.c"
Cohesion: 0.50
Nodes (6): main(), peek_buffer(), run_tracer(), run_victim(), touch(), wait_for_file()

### Community 186 - "Binary trust"
Cohesion: 0.25
Nodes (8): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, The dynamic loader, trust

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 195 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "promptPassword"
Cohesion: 0.18
Nodes (12): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), promptPassword(), promptPasswordWith() (+4 more)

### Community 199 - "sources_test.go"
Cohesion: 0.20
Nodes (9): SystemTrusted(), systemTool(), archiveMembers(), decompressed(), newSources(), tarMembers(), TestArchiveMembersIgnoresNonArchives(), TestSourcesHoldFilesAndArchiveMembers() (+1 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "run.go"
Cohesion: 0.36
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 202 - "DaemonEvent"
Cohesion: 0.10
Nodes (32): UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemon(), runDaemonUI() (+24 more)

### Community 203 - "Existing"
Cohesion: 0.18
Nodes (8): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), Existing, ExistingRule

### Community 204 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 205 - "find"
Cohesion: 0.22
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.14
Nodes (14): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+6 more)

### Community 209 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (46): recordTickCmd(), recordModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts() (+38 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 213 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 214 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 216 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 217 - "linkView"
Cohesion: 0.33
Nodes (5): linkTarget(), newLinkView(), entryKind, libPlan, linkView

### Community 218 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 219 - "SanitizeText"
Cohesion: 0.12
Nodes (12): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), scriptIdentity(), startTrustGuard(), trustStartupError(), binaryVetter, vetDecision (+4 more)

### Community 220 - "shQuote"
Cohesion: 0.19
Nodes (4): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite

### Community 222 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 223 - "compileGlobNames"
Cohesion: 0.15
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 224 - "systemd.go"
Cohesion: 0.07
Nodes (55): anyDaemonRunning(), offerDaemonStart(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markInstallComplete() (+47 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "sources"
Cohesion: 0.31
Nodes (4): pathByte(), pathRecord(), TestPathRecord(), sources

### Community 227 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 228 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), SetUserHomes(), TestOpenSystemPlacedRefusesUserHome()

### Community 229 - "logLine"
Cohesion: 0.14
Nodes (10): logHook, logLine, journalLine(), TestJournalLine(), logRow(), logStyle(), formatEntry(), sortedKeys() (+2 more)

### Community 232 - "procRec"
Cohesion: 0.40
Nodes (5): judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), procRec

### Community 234 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 235 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "GeneralToolsOf"
Cohesion: 0.50
Nodes (4): GeneralToolsOf(), namedGeneralTool(), ResolvesToGeneralTool(), TestNamedGeneralTool()

### Community 238 - "removeKeyAndEmptyDir"
Cohesion: 0.67
Nodes (3): removeKeyAndEmptyDir(), removeMasterKey(), TestRemoveKeyAndEmptyDir()

### Community 244 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 644 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **45 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `engine`, `ConfinedEntry`, `sync.Mutex`, `guardUnitTest`, `time.Time`, `DaemonEvent`, `canonicalBinaryPath`, `installPage`, `os.File`, `GuardEvent`, `TrustGuard`, `Resource`, `buildOneGuard`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.11921633250938386 - nodes in this community are weakly interconnected._
- **Why does `engine` connect `engine` to `go_pkg_strings`, `sync.Mutex`, `buildOneGuard`, `bpfstats/main.go`, `os.File`, `Guard`, `github.com/cilium/ebpf.Map`, `NetGuard`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.021455938697318006 - nodes in this community are weakly interconnected._
- **Why does `Monitor` connect `Monitor` to `go_pkg_strings`, `.populateWatchInodes`, `sync.Mutex`, `pathCache`, `monitorUnitTest`, `mountTable`, `.eventTypeAllowed`, `Resource`, `.filterEvent`, `mountsUnitTest`, `NetGuard`, `NetEventType`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._